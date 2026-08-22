import Sanscript from '@indic-transliteration/sanscript';

const input = document.getElementById('analyzeInput');
const inScript = document.getElementById('inputScript');
const outScript = document.getElementById('outputScript');
const container = document.getElementById('resultsContainer');
const rawJson = document.getElementById('rawJson');
const jsonBtn = document.getElementById('jsonBtn');

document.getElementById('btnAnalyzeTab').addEventListener('click', (e) => switchTab('analyzeTab', e.target));
document.getElementById('btnDhatuTab').addEventListener('click', (e) => switchTab('dhatuTab', e.target));
document.getElementById('btnGenVerbTab').addEventListener('click', (e) => switchTab('genVerbTab', e.target));
document.getElementById('btnGenPartTab').addEventListener('click', (e) => switchTab('genPartTab', e.target));
document.getElementById('btnGenNounTab').addEventListener('click', (e) => switchTab('genNounTab', e.target));

function switchTab(tabId, btnTarget) {
    document.querySelectorAll('.tab-content').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('.tab-btn').forEach(el => el.classList.remove('active'));
    document.getElementById(tabId).classList.add('active');
    btnTarget.classList.add('active');
    container.innerHTML = ''; jsonBtn.style.display = 'none'; rawJson.style.display = 'none';
}

jsonBtn.addEventListener('click', () => { rawJson.style.display = rawJson.style.display === 'none' ? 'block' : 'none'; });

document.querySelectorAll('input').forEach(inp => {
    inp.addEventListener('input', () => { if (/[\u0900-\u097F]/.test(inp.value)) inScript.value = 'devanagari'; });
});

function norm(word) {
    if (!word) return word;
    if (/[\u0900-\u097F]/.test(word)) { inScript.value = 'devanagari'; return Sanscript.t(word, 'devanagari', 'slp1'); }
    if (/[āīūṛṝḷḹēōṃḥśṣṭḍṇñṅ]/.test(word.toLowerCase())) { inScript.value = 'iast'; return Sanscript.t(word, 'iast', 'slp1'); }
    return inScript.value === 'slp1' ? word : Sanscript.t(word, inScript.value, 'slp1');
}

function out(text, col) {
    if (!text || text === '-') return '-';
    if (['Dhatu ID'].includes(col)) return text;
    if (/^[0-9.]+$/.test(text)) return text;
    
    if (['Meaning', 'Prefixed Meaning'].includes(col) || /[\u0900-\u097F]/.test(text)) return text;

    const rawTags = [
        'masc/fem', 'masc/neut', 'masc/fem/neut', 'any',
        'mula', 'base'
    ];
    return text.split(' ').map(w => rawTags.includes(w.toLowerCase()) ? w : Sanscript.t(w, 'slp1', outScript.value)).join(' ');
}

function buildTable(title, cols, rows) {
    if (!rows || rows.length === 0) return '';
    let html = `<h2 class="category-title">${title}</h2><table><thead><tr>`;
    cols.forEach(c => html += `<th>${c}</th>`);
    html += `</tr></thead><tbody>`;
    rows.forEach(r => {
        html += `<tr>`;
        cols.forEach(c => html += `<td>${out(r[c.toLowerCase().replace(/ /g, '_')] || r[c] || '-', c)}</td>`);
        html += `</tr>`;
    });
    return html + `</tbody></table>`;
}

function buildDeclensionGrid(title, declensions) {
    const cases = ['prathama', 'dvitiya', 'tritiya', 'caturthi', 'panchami', 'sasthi', 'saptami', 'sambodhana'];
    let html = `<h2 class="category-title">${title}</h2><table><thead><tr><th>Case</th><th>Eka</th><th>Dvi</th><th>Bahu</th></tr></thead><tbody>`;
    cases.forEach(c => {
        if(declensions[c]) {
            html += `<tr><td>${c.charAt(0).toUpperCase() + c.slice(1)}</td><td>${out(declensions[c][0])}</td><td>${out(declensions[c][1])}</td><td>${out(declensions[c][2])}</td></tr>`;
        }
    });
    return html + `</tbody></table>`;
}

async function fetchAPI(url) {
    container.innerHTML = "<p style='text-align:center;'>Querying Database...</p>";
    document.body.classList.add('loading');
    jsonBtn.style.display = 'none'; rawJson.style.display = 'none';
    try {
        const response = await fetch(url);
        const text = await response.text();
        if (!response.ok) {
            if (response.status === 503 || response.status === 429) {
                throw new Error(`Server waking up (HTTP ${response.status}) — Render free tier hibernates after inactivity. Please wait 20s and retry.`);
            }
            throw new Error(`HTTP ${response.status}: ${text.slice(0,200) || response.statusText}`);
        }
        if (!text) {
            throw new Error("Empty response from server — please retry (server may be waking up)");
        }
        let data;
        try { data = JSON.parse(text); } catch (e) {
            throw new Error(`Invalid JSON response: ${text.slice(0,200)}`);
        }
        rawJson.innerText = JSON.stringify(data, null, 2);
        return data;
    } catch (err) {
        container.innerHTML = `<div class='error-msg'>Server Error: ${err.message}</div>`;
        return null;
    } finally { document.body.classList.remove('loading'); }
}

async function runAnalyzer() {
    const rawWord = document.getElementById('analyzeInput').value.trim();
    if (!rawWord) return;
    const slp1Word = norm(rawWord);
    
    const data = await fetchAPI(`/api/analyze/${slp1Word}`);
    if (!data) return;

    let html = '';
    html += buildTable('Verbs (Tiṅanta)', ['Dhatu ID', 'Root', 'Meaning', 'Prefixed Meaning', 'Upasarga', 'Lakara', 'Purusha', 'Vacana', 'Voice'], data.verbs);
    html += buildTable('Declensions (Subanta)', ['Base Form', 'Dhatu ID', 'Upasarga', 'Prefixed Meaning', 'Pratyaya', 'Gender', 'Case', 'Vacana'], data.declensions);
    html += buildTable('Participles / Avyayas', ['Base Form', 'Root', 'Pratyaya', 'Dhatu ID', 'Upasarga', 'Prefixed Meaning', 'Gender', 'Case', 'Vacana'], data.participles);
    html += buildTable('Pronouns', ['Base Form', 'Gender', 'Case', 'Vacana'], data.pronouns);
    html += buildTable('Numerals', ['Base Form', 'Gender', 'Case', 'Vacana'], data.numerals);
    html += buildTable('Irregular Nouns', ['Base Form', 'Gender', 'Case', 'Vacana'], data.irregulars);

    if (!html) {
        html = `<div class='error-msg'>No matches found for "${rawWord}" (parsed as SLP1: ${slp1Word}).<br>Try entering an un-sandhied word.</div>`;
    }

    container.innerHTML = html;
    rawJson.innerText = JSON.stringify(data, null, 2);
    if(html && !html.includes('error-msg')) jsonBtn.style.display = 'block';
}

async function runDhatuSearch() {
    const rawWord = document.getElementById('dhatuInput').value.trim();
    if (!rawWord) return;
    const query = norm(rawWord);
    const data = await fetchAPI(`/api/dhatus/${query}`);
    if (!data || data.length === 0) { container.innerHTML = `<div class='error-msg'>No Dhātus found matching "${rawWord}".</div>`; return; }
    
    data.sort((a,b) => (a.dhatu_id || '').localeCompare(b.dhatu_id || ''));
    container.innerHTML = buildTable('Dhātu Results', ['Dhatu ID', 'Root', 'Gana', 'Meaning', 'Upasarga', 'Pratyaya', 'Base Form'], data);
    jsonBtn.style.display = 'block';
}

document.getElementById('runAnalyzerBtn').addEventListener('click', runAnalyzer);
document.getElementById('runDhatuBtn').addEventListener('click', runDhatuSearch);
document.getElementById('runGenVerbBtn').addEventListener('click', async () => {
    const r = norm(document.getElementById('gvRoot').value.trim());
    const p = new URLSearchParams({ root: r, upasarga: norm(document.getElementById('gvUpa').value.trim()), lakara: document.getElementById('gvLakara').value, purusha: document.getElementById('gvPurusha').value, voice: document.getElementById('gvVoice').value, prayoga: document.getElementById('gvPrayoga').value, derivative: document.getElementById('gvDerivative').value });
    if (!r) return alert("Please enter a root.");
    const data = await fetchAPI(`/api/generate/verb?${p}`);
    if (data.error) { container.innerHTML = `<div class='error-msg'>${data.error}</div>`; return; }
    container.innerHTML = buildTable('Generated Verb', ['Eka', 'Dvi', 'Bahu'], [data]);
    jsonBtn.style.display = 'block';
});

document.getElementById('runGenPartBtn').addEventListener('click', async () => {
    const r = norm(document.getElementById('gpRoot').value.trim());
    const p = new URLSearchParams({ root: r, upasarga: norm(document.getElementById('gpUpa').value.trim()), pratyaya: document.getElementById('gpPratyaya').value, gender: document.getElementById('gpGender').value, derivative: document.getElementById('gpDerivative').value });
    if (!r) return alert("Please enter a root.");
    const data = await fetchAPI(`/api/generate/participle?${p}`);
    if (data.error) { container.innerHTML = `<div class='error-msg'>${data.error}</div>`; return; }
    if (data.type === 'avyaya') { container.innerHTML = buildTable('Generated Avyaya', ['Base Form'], [data]); } 
    else { container.innerHTML = buildDeclensionGrid(`Declension: ${out(data.base_form)}`, data.declensions); }
    jsonBtn.style.display = 'block';
});

document.getElementById('runGenNounBtn').addEventListener('click', async () => {
    const b = norm(document.getElementById('gnBase').value.trim());
    const p = new URLSearchParams({ base: b, gender: document.getElementById('gnGender').value });
    if (!b) return alert("Please enter a base noun.");
    const data = await fetchAPI(`/api/generate/declension?${p}`);
    if (data.error) { container.innerHTML = `<div class='error-msg'>${data.error}</div>`; return; }
    container.innerHTML = buildDeclensionGrid(`Declension: ${out(b)}`, data.declensions);
    jsonBtn.style.display = 'block';
});

document.getElementById('analyzeInput').addEventListener('keypress', (e) => { if (e.key === 'Enter') runAnalyzer(); });
document.getElementById('dhatuInput').addEventListener('keypress', (e) => { if (e.key === 'Enter') runDhatuSearch(); });
outScript.addEventListener('change', () => { 
    const active = document.querySelector('.tab-btn.active').id;
    if (container.innerHTML && !container.innerHTML.includes('No matches') && !container.innerHTML.includes('Error')) {
        if (active === 'btnAnalyzeTab') runAnalyzer();
        if (active === 'btnDhatuTab') runDhatuSearch();
    }
});
