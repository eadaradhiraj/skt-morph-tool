import Sanscript from '@indic-transliteration/sanscript';

const input = document.getElementById('analyzeInput');
const genRoot = document.getElementById('genRoot');
const inScript = document.getElementById('inputScript');
const outScript = document.getElementById('outputScript');
const container = document.getElementById('resultsContainer');
const rawJson = document.getElementById('rawJson');
const jsonBtn = document.getElementById('jsonBtn');

document.getElementById('btnAnalyzeTab').addEventListener('click', (e) => switchTab('analyzeTab', e.target));
document.getElementById('btnGenerateTab').addEventListener('click', (e) => switchTab('generateTab', e.target));

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

function normalizeToSLP1(word) {
    if (!word) return word;
    if (inScript.value === 'slp1') return word;
    return Sanscript.t(word, inScript.value, 'slp1');
}

function renderOutput(slp1Text, columnTitle) {
    if (!slp1Text || slp1Text === '-') return '-';
    
    if (['Dhatu ID', 'Meaning Hint'].includes(columnTitle)) return slp1Text;
    if (/^[0-9.]+$/.test(slp1Text)) return slp1Text;

    const rawTags = ['masc', 'fem', 'neut', 'eka', 'dvi', 'bahu', 'masc/fem', 'masc/neut', 'masc/fem/neut', 'any'];
    
    return slp1Text.split(' ').map(word => {
        if (rawTags.includes(word.toLowerCase())) return word;
        return Sanscript.t(word, 'slp1', outScript.value);
    }).join(' ');
}

function createTable(title, columns, dataRows) {
    if (!dataRows || dataRows.length === 0) return '';
    let html = `<h2 class="category-title">${title}</h2><table><thead><tr>`;
    columns.forEach(c => html += `<th>${c}</th>`);
    html += `</tr></thead><tbody>`;
    dataRows.forEach(row => {
        html += `<tr>`;
        columns.forEach(c => {
            const key = c.toLowerCase().replace(' ', '_');
            const val = row[key] || row[c] || '-';
            html += `<td>${renderOutput(val, c)}</td>`;
        });
        html += `</tr>`;
    });
    html += `</tbody></table>`;
    return html;
}

async function fetchAPI(url) {
    container.innerHTML = "<p style='text-align:center;'>Querying Database...</p>";
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
    } finally {
        document.body.classList.remove('loading');
    }
}

async function runAnalyzer() {
    const rawWord = input.value.trim();
    if (!rawWord) return;
    const slp1Word = normalizeToSLP1(rawWord);
    
    const data = await fetchAPI(`/api/analyze/${slp1Word}`);
    if (!data) return;

    let html = '';
    html += createTable('Verbs (Tiṅanta)', ['Dhatu ID', 'Upasarga', 'Lakara', 'Purusha', 'Vacana', 'Voice'], data.verbs);
    html += createTable('Declensions (Subanta)', ['Base Form', 'Dhatu ID', 'Upasarga', 'Pratyaya', 'Gender', 'Case', 'Vacana'], data.declensions);
    html += createTable('Participles / Avyayas', ['Base Form', 'Pratyaya', 'Dhatu ID', 'Upasarga'], data.participles);
    html += createTable('Namadhatus', ['Base Noun', 'Pratyaya', 'Meaning Hint'], data.namadhatus);
    html += createTable('Pronouns', ['Base Form', 'Gender', 'Case', 'Vacana'], data.pronouns);
    html += createTable('Numerals', ['Base Form', 'Gender', 'Case', 'Vacana'], data.numerals);
    html += createTable('Irregular Nouns', ['Base Form', 'Gender', 'Case', 'Vacana'], data.irregulars);

    if (!html) html = `<div class='error-msg'>No matches found for "${rawWord}". Try an un-sandhied word.</div>`;
    container.innerHTML = html;
    if(html && !html.includes('error-msg')) jsonBtn.style.display = 'block';
}

async function runGenerator() {
    const rawRoot = genRoot.value.trim();
    const rawUpa = document.getElementById('genUpa').value.trim();
    if (!rawRoot) return alert("Please enter a root.");

    const root = normalizeToSLP1(rawRoot);
    const upa = normalizeToSLP1(rawUpa);

    const params = new URLSearchParams({
        root: root, upasarga: upa,
        lakara: document.getElementById('genLakara').value,
        purusha: document.getElementById('genPurusha').value,
        voice: document.getElementById('genVoice').value
    });

    const data = await fetchAPI(`/api/generate/verb?${params}`);
    if (!data || data.error) {
        container.innerHTML = `<div class='error-msg'>${data ? data.error : "Generation Failed"}</div>`;
        return;
    }

    container.innerHTML = createTable('Generated Verb', ['Eka', 'Dvi', 'Bahu'], [data]);
    jsonBtn.style.display = 'block';
}

document.getElementById('btnAnalyze').addEventListener('click', runAnalyzer);
document.getElementById('btnGenerate').addEventListener('click', runGenerator);
input.addEventListener('keypress', (e) => { if (e.key === 'Enter') runAnalyzer(); });
genRoot.addEventListener('keypress', (e) => { if (e.key === 'Enter') runGenerator(); });

outScript.addEventListener('change', () => { 
    const activeTab = document.querySelector('.tab-btn.active').innerText;
    if (container.innerHTML && !container.innerHTML.includes('No matches') && !container.innerHTML.includes('Error')) {
        if (activeTab === 'Analyzer') runAnalyzer();
        if (activeTab === 'Generator') runGenerator();
    }
});
