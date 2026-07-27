import Sanscript from '@indic-transliteration/sanscript';

const inScript = document.getElementById('inputScript');
const outScript = document.getElementById('outputScript');
const container = document.getElementById('resultsContainer');
const rawJson = document.getElementById('rawJson');
const jsonBtn = document.getElementById('jsonBtn');

document.getElementById('btnAnalyzeTab').addEventListener('click', (e) => switchTab('analyzeTab', e.target));
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
    return inScript.value === 'slp1' ? word : Sanscript.t(word, inScript.value, 'slp1');
}

function out(text, col) {
    if (!text || text === '-') return '-';
    if (['Dhatu ID'].includes(col)) return text;
    if (/^[0-9.]+$/.test(text)) return text;
    const rawTags = ['masc', 'fem', 'neut', 'eka', 'dvi', 'bahu', 'masc/fem', 'masc/neut', 'masc/fem/neut', 'any'];
    return text.split(' ').map(w => rawTags.includes(w.toLowerCase()) ? w : Sanscript.t(w, 'slp1', outScript.value)).join(' ');
}

function buildTable(title, cols, rows) {
    if (!rows || rows.length === 0) return '';
    let html = `<h2 class="category-title">${title}</h2><table><thead><tr>`;
    cols.forEach(c => html += `<th>${c}</th>`);
    html += `</tr></thead><tbody>`;
    rows.forEach(r => {
        html += `<tr>`;
        cols.forEach(c => html += `<td>${out(r[c.toLowerCase().replace(' ', '_')] || r[c] || '-', c)}</td>`);
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
    container.innerHTML = "<p style='text-align:center;'>Querying...</p>";
    document.body.classList.add('loading');
    jsonBtn.style.display = 'none'; rawJson.style.display = 'none';
    try {
        const response = await fetch(url);
        const data = await response.json();
        rawJson.innerText = JSON.stringify(data, null, 2);
        return data;
    } catch (err) {
        container.innerHTML = `<div class='error-msg'>Server Error: ${err.message}</div>`;
        throw err;
    } finally { document.body.classList.remove('loading'); }
}

document.getElementById('runAnalyzerBtn').addEventListener('click', async () => {
    const rawWord = document.getElementById('analyzeInput').value.trim();
    if (!rawWord) return;
    const data = await fetchAPI(`/api/analyze/${norm(rawWord)}`);
    if (!data) return;
    let html = buildTable('Verbs (Tiṅanta)', ['Dhatu ID', 'Upasarga', 'Lakara', 'Purusha', 'Vacana', 'Voice'], data.verbs);
    html += buildTable('Declensions (Subanta)', ['Base Form', 'Dhatu ID', 'Upasarga', 'Pratyaya', 'Gender', 'Case', 'Vacana'], data.declensions);
    html += buildTable('Participles / Avyayas', ['Base Form', 'Pratyaya', 'Dhatu ID', 'Upasarga'], data.participles);
    html += buildTable('Pronouns', ['Base Form', 'Gender', 'Case', 'Vacana'], data.pronouns);
    html += buildTable('Numerals', ['Base Form', 'Gender', 'Case', 'Vacana'], data.numerals);
    html += buildTable('Irregular Nouns', ['Base Form', 'Gender', 'Case', 'Vacana'], data.irregulars);
    container.innerHTML = html || `<div class='error-msg'>No matches found.</div>`;
    if(html) jsonBtn.style.display = 'block';
});

document.getElementById('runGenVerbBtn').addEventListener('click', async () => {
    const r = norm(document.getElementById('gvRoot').value.trim());
    if (!r) return alert("Please enter a root.");
    const p = new URLSearchParams({ 
        root: r, 
        upasarga: norm(document.getElementById('gvUpa').value.trim()), 
        lakara: document.getElementById('gvLakara').value, 
        purusha: document.getElementById('gvPurusha').value, 
        voice: document.getElementById('gvVoice').value,
        prayoga: document.getElementById('gvPrayoga').value,
        derivative: document.getElementById('gvDerivative').value
    });
    const data = await fetchAPI(`/api/generate/verb?${p}`);
    if (data.error) { container.innerHTML = `<div class='error-msg'>${data.error}</div>`; return; }
    container.innerHTML = buildTable('Generated Verb', ['Eka', 'Dvi', 'Bahu'], [data]);
    jsonBtn.style.display = 'block';
});

document.getElementById('runGenPartBtn').addEventListener('click', async () => {
    const r = norm(document.getElementById('gpRoot').value.trim());
    if (!r) return alert("Please enter a root.");
    const p = new URLSearchParams({ 
        root: r, 
        upasarga: norm(document.getElementById('gpUpa').value.trim()), 
        pratyaya: document.getElementById('gpPratyaya').value, 
        gender: document.getElementById('gpGender').value,
        derivative: document.getElementById('gpDerivative').value 
    });
    const data = await fetchAPI(`/api/generate/participle?${p}`);
    if (data.error) { container.innerHTML = `<div class='error-msg'>${data.error}</div>`; return; }
    if (data.type === 'avyaya') { container.innerHTML = buildTable('Generated Avyaya', ['Base Form'], [data]); } 
    else { container.innerHTML = buildDeclensionGrid(`Declension: ${out(data.base_form)}`, data.declensions); }
    jsonBtn.style.display = 'block';
});

document.getElementById('runGenNounBtn').addEventListener('click', async () => {
    const b = norm(document.getElementById('gnBase').value.trim());
    if (!b) return alert("Please enter a base noun.");
    const p = new URLSearchParams({ base: b, gender: document.getElementById('gnGender').value });
    const data = await fetchAPI(`/api/generate/declension?${p}`);
    if (data.error) { container.innerHTML = `<div class='error-msg'>${data.error}</div>`; return; }
    container.innerHTML = buildDeclensionGrid(`Declension: ${out(b)}`, data.declensions);
    jsonBtn.style.display = 'block';
});

// Event listeners for enter key
document.getElementById('analyzeInput').addEventListener('keypress', (e) => { if (e.key === 'Enter') document.getElementById('runAnalyzerBtn').click(); });
document.getElementById('gvRoot').addEventListener('keypress', (e) => { if (e.key === 'Enter') document.getElementById('runGenVerbBtn').click(); });
document.getElementById('gpRoot').addEventListener('keypress', (e) => { if (e.key === 'Enter') document.getElementById('runGenPartBtn').click(); });
document.getElementById('gnBase').addEventListener('keypress', (e) => { if (e.key === 'Enter') document.getElementById('runGenNounBtn').click(); });

outScript.addEventListener('change', () => { 
    const active = document.querySelector('.tab-btn.active').id;
    if (container.innerHTML && !container.innerHTML.includes('No matches') && !container.innerHTML.includes('Error')) {
        if (active === 'btnAnalyzeTab') document.getElementById('runAnalyzerBtn').click();
        if (active === 'btnGenVerbTab') document.getElementById('runGenVerbBtn').click();
        if (active === 'btnGenPartTab') document.getElementById('runGenPartBtn').click();
        if (active === 'btnGenNounTab') document.getElementById('runGenNounBtn').click();
    }
});
