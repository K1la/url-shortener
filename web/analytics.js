const api = {
  analytics: async (alias) => {
    const res = await fetch(`/api/analytics/${encodeURIComponent(alias)}`);
    if (!res.ok) {
      const text = await res.text().catch(() => '');
      throw new Error(text || 'Failed to fetch analytics');
    }
    return res.json();
  },
};

function byId(id) { return document.getElementById(id); }
function show(el) { el.classList.remove('hidden'); }
function hide(el) { el.classList.add('hidden'); }

// Analytics form
const analyticsForm = byId('analytics-form');
const aliasInput = byId('alias-input');
const analyticsResults = byId('analytics-results');
const analyticsError = byId('analytics-error');
const loading = byId('loading');
const emptyState = byId('empty-state');

// Results elements
const totalClicksEl = byId('total-clicks');
const uniqueUaEl = byId('unique-ua');
const todayClicksEl = byId('today-clicks');
const tableBody = byId('analytics-table-body');
const chartPlaceholder = byId('chart-placeholder');

// Period tabs
const periodTabs = document.querySelectorAll('.period-tab');
let currentData = null;
let currentPeriod = 'daily';

// Form submission
analyticsForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  hide(analyticsError);
  hide(analyticsResults);
  hide(emptyState);
  show(loading);
  
  try {
    const alias = aliasInput.value.trim();
    if (!alias) return;
    
    const data = await api.analytics(alias);
    const result = data && (data.result || data.Result || data);
    
    if (!result || result.total_clicks === 0) {
      hide(loading);
      show(emptyState);
      return;
    }
    
    currentData = result;
    displayAnalytics(result);
    
    hide(loading);
    show(analyticsResults);
  } catch (err) {
    hide(loading);
    analyticsError.textContent = err.message || 'Ошибка при получении аналитики';
    show(analyticsError);
  }
});

// Period tabs
periodTabs.forEach(tab => {
  tab.addEventListener('click', () => {
    periodTabs.forEach(t => t.classList.remove('active'));
    tab.classList.add('active');
    currentPeriod = tab.dataset.period;
    
    if (currentData) {
      displayAnalytics(currentData);
    }
  });
});

function displayAnalytics(data) {
  // Update stats
  const total = data.total_clicks || 0;
  const daily = data.daily || {};
  const monthly = data.monthly || {};
  const userAgent = data.user_agent || {};
  
  totalClicksEl.textContent = total.toString();
  uniqueUaEl.textContent = Object.keys(userAgent).length.toString();
  
  // Today's clicks
  const today = new Date().toISOString().split('T')[0];
  const todayClicks = daily[today] || 0;
  todayClicksEl.textContent = todayClicks.toString();
  
  // Update table
  updateTable(currentPeriod === 'daily' ? daily : monthly);
  
  // Update chart placeholder
  updateChart(currentPeriod === 'daily' ? daily : monthly);
}

function updateTable(data) {
  tableBody.innerHTML = '';
  
  const entries = Object.entries(data).sort((a, b) => {
    // Sort by date descending
    return currentPeriod === 'daily' ? b[0].localeCompare(a[0]) : b[0].localeCompare(a[0]);
  });
  
  const total = Object.values(data).reduce((sum, count) => sum + count, 0);
  
  entries.forEach(([period, clicks]) => {
    const percentage = total > 0 ? ((clicks / total) * 100).toFixed(1) : 0;
    
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td>${formatPeriod(period)}</td>
      <td>${clicks}</td>
      <td>${percentage}%</td>
    `;
    tableBody.appendChild(tr);
  });
}

function updateChart(data) {
  const entries = Object.entries(data).sort((a, b) => a[0].localeCompare(b[0]));
  
  if (entries.length === 0) {
    chartPlaceholder.innerHTML = '<div style="text-align: center; color: var(--muted);">📊 Нет данных для отображения</div>';
    return;
  }
  
  // Create a simple HTML chart with bars
  const maxClicks = Math.max(...Object.values(data));
  const chartHTML = entries.map(([period, clicks]) => {
    const percentage = maxClicks > 0 ? (clicks / maxClicks) * 100 : 0;
    return `
      <div style="display: flex; align-items: center; margin-bottom: 0.5rem; padding: 0.5rem; background: var(--card); border-radius: 6px; border: 1px solid var(--border);">
        <div style="min-width: 80px; font-weight: 600; color: var(--text);">${formatPeriod(period)}</div>
        <div style="flex: 1; margin: 0 1rem; position: relative;">
          <div style="background: var(--border); height: 8px; border-radius: 4px; overflow: hidden;">
            <div style="background: linear-gradient(90deg, var(--accent), var(--accent-2)); height: 100%; width: ${percentage}%; transition: width 0.3s ease;"></div>
          </div>
        </div>
        <div style="min-width: 40px; text-align: right; font-weight: 600; color: var(--text);">${clicks}</div>
      </div>
    `;
  }).join('');
  
  chartPlaceholder.innerHTML = chartHTML;
}

function formatPeriod(period) {
  if (currentPeriod === 'daily') {
    // Format date as DD.MM.YYYY
    const [year, month, day] = period.split('-');
    return `${day}.${month}.${year}`;
  } else {
    // Format month as MM.YYYY
    const [year, month] = period.split('-');
    return `${month}.${year}`;
  }
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
  // Focus on input
  aliasInput.focus();
});
