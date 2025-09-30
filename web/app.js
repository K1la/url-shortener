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
    // API returns { result: { url, short_url, ... } }
    const result = data && (data.result || data.Result || data);
    const alias = (result && (result.short_url || result.ShortURL)) || payload.user_short_url;
    const base = window.location.origin;
    const href = alias ? `${base}/api/s/${alias}` : (result.url || result.URL || base);
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

// Open short link
const openForm = byId('open-form');
const openAlias = byId('open-alias');
openForm.addEventListener('submit', (e) => {
  e.preventDefault();
  const alias = openAlias.value.trim();
  if (!alias) return;
  const href = `${window.location.origin}/api/s/${encodeURIComponent(alias)}`;
  window.open(href, '_blank');
});

// Analytics removed - now handled on separate page

function escapeHtml(str){
  return String(str || '')
    .replace(/&/g,'&amp;')
    .replace(/</g,'&lt;')
    .replace(/>/g,'&gt;')
    .replace(/"/g,'&quot;')
    .replace(/'/g,'&#039;');
}


