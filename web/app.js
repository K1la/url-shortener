const api = {
  create: async ({ url, user_short_url }) => {
    const res = await fetch('/api/shorten', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url, user_short_url }),
    });
    if (!res.ok) {
      const text = await res.text().catch(() => '');
      throw new Error(text || 'Failed to create short link');
    }
    return res.json().catch(() => ({}));
  },
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

// Create link
const createForm = byId('create-form');
const longUrlInput = byId('long-url');
const customAliasInput = byId('custom-alias');
const resultBox = byId('result');
const shortLink = byId('short-link');
const copyBtn = byId('copy-btn');
const createErr = byId('create-error');

createForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  hide(createErr);
  hide(resultBox);
  try {
    const payload = {
      url: longUrlInput.value.trim(),
      user_short_url: customAliasInput.value.trim() || undefined,
    };
    const data = await api.create(payload);
    // Try to detect returned fields; gracefully build link
    const alias = data.short_url || data.ShortURL || payload.user_short_url;
    const base = window.location.origin;
    const href = alias ? `${base}/s/${alias}` : (data.url || data.URL || base);
    shortLink.textContent = href;
    shortLink.href = href;
    show(resultBox);
  } catch (err) {
    createErr.textContent = err.message || 'Ошибка при создании ссылки';
    show(createErr);
  }
});

copyBtn.addEventListener('click', async () => {
  const text = shortLink.href;
  try {
    await navigator.clipboard.writeText(text);
    copyBtn.textContent = 'Скопировано';
    setTimeout(() => (copyBtn.textContent = 'Копировать'), 1200);
  } catch (_) {
    // noop
  }
});

// Analytics
const analyticsForm = byId('analytics-form');
const aliasInput = byId('alias-for-analytics');
const analyticsBox = byId('analytics');
const analyticsErr = byId('analytics-error');
const totalClicksEl = byId('total-clicks');
const uniqueUaEl = byId('unique-ua');
const clicksBody = byId('clicks-body');

analyticsForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  hide(analyticsErr);
  hide(analyticsBox);
  clicksBody.innerHTML = '';
  try {
    const alias = aliasInput.value.trim();
    const data = await api.analytics(alias);
    // Expecting structure like: { total, clicks: [{time, device, os, browser, ip, user_agent}], ... }
    const clicks = data.clicks || data.items || data || [];
    totalClicksEl.textContent = (data.total ?? clicks.length).toString();

    const uas = new Set();
    clicks.forEach((c) => {
      const when = c.time || c.created_at || c.CreatedAt || c.timestamp;
      const device = c.device || c.Device || '';
      const os = c.os || c.OS || '';
      const browser = c.browser || c.Browser || '';
      const ip = c.ip || c.IP || '';
      const ua = c.user_agent || c.UserAgent || '';
      if (ua) uas.add(ua);
      const dt = when ? new Date(when) : null;
      const dtStr = dt ? dt.toLocaleString() : '';
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td>${dtStr}</td>
        <td>${escapeHtml(device)}</td>
        <td>${escapeHtml(os)}</td>
        <td>${escapeHtml(browser)}</td>
        <td>${escapeHtml(ip)}</td>
        <td>${escapeHtml(ua)}</td>
      `;
      clicksBody.appendChild(tr);
    });
    uniqueUaEl.textContent = uas.size.toString();
    show(analyticsBox);
  } catch (err) {
    analyticsErr.textContent = err.message || 'Ошибка при получении аналитики';
    show(analyticsErr);
  }
});

function escapeHtml(str){
  return String(str || '')
    .replace(/&/g,'&amp;')
    .replace(/</g,'&lt;')
    .replace(/>/g,'&gt;')
    .replace(/"/g,'&quot;')
    .replace(/'/g,'&#039;');
}


