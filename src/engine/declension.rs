use fancy_regex::Regex;
use std::collections::HashMap;

pub fn apply_natva(word: &str) -> String {
    let pattern = Regex::new(r"([rfFz][aAiIuUfFxXeEoOkKgGNpPbBmyvhM]*)n(?=[aAiIuUfFxXeEoOmyv])").unwrap();
    let mut current = word.to_string();
    while pattern.is_match(&current).unwrap_or(false) {
        current = pattern.replace(&current, "${1}R").into_owned();
    }
    current
}

fn build_grid(pr: [&str;3], dv: [&str;3], tr: [&str;3], ca: [&str;3], pa: [&str;3], sa: [&str;3], sap: [&str;3], sam: [&str;3]) -> HashMap<String, Vec<String>> {
    let mut map = HashMap::new();
    map.insert("prathama".to_string(), vec![pr[0].to_string(), pr[1].to_string(), pr[2].to_string()]);
    map.insert("dvitiya".to_string(), vec![dv[0].to_string(), dv[1].to_string(), dv[2].to_string()]);
    map.insert("tritiya".to_string(), vec![tr[0].to_string(), tr[1].to_string(), tr[2].to_string()]);
    map.insert("caturthi".to_string(), vec![ca[0].to_string(), ca[1].to_string(), ca[2].to_string()]);
    map.insert("panchami".to_string(), vec![pa[0].to_string(), pa[1].to_string(), pa[2].to_string()]);
    map.insert("sasthi".to_string(), vec![sa[0].to_string(), sa[1].to_string(), sa[2].to_string()]);
    map.insert("saptami".to_string(), vec![sap[0].to_string(), sap[1].to_string(), sap[2].to_string()]);
    map.insert("sambodhana".to_string(), vec![sam[0].to_string(), sam[1].to_string(), sam[2].to_string()]);
    map
}

pub fn decline_noun(base: &str, gender: &str) -> Result<HashMap<String, Vec<String>>, String> {
    if let Some(irreg) = crate::engine::irregulars::decline_irregular(base, gender) {
        let mut fixed = HashMap::new();
        for (k, v) in irreg {
            let natva_v: Vec<String> = v.into_iter().map(|w| apply_natva(&w)).collect();
            fixed.insert(k, natva_v);
        }
        return Ok(fixed);
    }

    let mut res = HashMap::new();
    let s = if base.len() > 0 { &base[..base.len()-1] } else { "" };
    
    if base.ends_with('a') {
        if gender == "masculine" { res = build_grid([&format!("{}aH",s), &format!("{}O",s), &format!("{}AH",s)], [&format!("{}am",s), &format!("{}O",s), &format!("{}An",s)], [&format!("{}ena",s), &format!("{}AByAm",s), &format!("{}EH",s)], [&format!("{}Aya",s), &format!("{}AByAm",s), &format!("{}eByaH",s)], [&format!("{}At",s), &format!("{}AByAm",s), &format!("{}eByaH",s)], [&format!("{}asya",s), &format!("{}ayoH",s), &format!("{}AnAm",s)], [&format!("{}e",s), &format!("{}ayoH",s), &format!("{}ezu",s)], [&format!("{}a",s), &format!("{}O",s), &format!("{}AH",s)]); }
        else if gender == "neuter" { res = build_grid([&format!("{}am",s), &format!("{}e",s), &format!("{}Ani",s)], [&format!("{}am",s), &format!("{}e",s), &format!("{}Ani",s)], [&format!("{}ena",s), &format!("{}AByAm",s), &format!("{}EH",s)], [&format!("{}Aya",s), &format!("{}AByAm",s), &format!("{}eByaH",s)], [&format!("{}At",s), &format!("{}AByAm",s), &format!("{}eByaH",s)], [&format!("{}asya",s), &format!("{}ayoH",s), &format!("{}AnAm",s)], [&format!("{}e",s), &format!("{}ayoH",s), &format!("{}ezu",s)], [&format!("{}a",s), &format!("{}e",s), &format!("{}Ani",s)]); }
    } else if base.ends_with('A') {
        if gender == "feminine" { res = build_grid([&format!("{}A",s), &format!("{}e",s), &format!("{}AH",s)], [&format!("{}Am",s), &format!("{}e",s), &format!("{}AH",s)], [&format!("{}ayA",s), &format!("{}AByAm",s), &format!("{}ABiH",s)], [&format!("{}AyE",s), &format!("{}AByAm",s), &format!("{}AByaH",s)], [&format!("{}AyAH",s), &format!("{}AByAm",s), &format!("{}AByaH",s)], [&format!("{}AyAH",s), &format!("{}ayoH",s), &format!("{}AnAm",s)], [&format!("{}AyAm",s), &format!("{}ayoH",s), &format!("{}Asu",s)], [&format!("{}e",s), &format!("{}e",s), &format!("{}AH",s)]); }
    } else if base.ends_with('I') {
        if gender == "feminine" { res = build_grid([&format!("{}I",s), &format!("{}yO",s), &format!("{}yaH",s)], [&format!("{}Im",s), &format!("{}yO",s), &format!("{}IH",s)], [&format!("{}yA",s), &format!("{}IByAm",s), &format!("{}IBiH",s)], [&format!("{}yE",s), &format!("{}IByAm",s), &format!("{}IByaH",s)], [&format!("{}yAH",s), &format!("{}IByAm",s), &format!("{}IByaH",s)], [&format!("{}yAH",s), &format!("{}yoH",s), &format!("{}InAm",s)], [&format!("{}yAm",s), &format!("{}yoH",s), &format!("{}Izu",s)], [&format!("{}i",s), &format!("{}yO",s), &format!("{}yaH",s)]); }
    } else if base.ends_with('i') {
        if gender == "masculine" { res = build_grid([&format!("{}iH",s), &format!("{}I",s), &format!("{}ayaH",s)], [&format!("{}im",s), &format!("{}I",s), &format!("{}In",s)], [&format!("{}inA",s), &format!("{}iByAm",s), &format!("{}iBiH",s)], [&format!("{}aye",s), &format!("{}iByAm",s), &format!("{}iByaH",s)], [&format!("{}eH",s), &format!("{}iByAm",s), &format!("{}iByaH",s)], [&format!("{}eH",s), &format!("{}yoH",s), &format!("{}InAm",s)], [&format!("{}O",s), &format!("{}yoH",s), &format!("{}izu",s)], [&format!("{}e",s), &format!("{}I",s), &format!("{}ayaH",s)]); }
        else if gender == "feminine" { res = build_grid([&format!("{}iH",s), &format!("{}I",s), &format!("{}ayaH",s)], [&format!("{}im",s), &format!("{}I",s), &format!("{}IH",s)], [&format!("{}yA",s), &format!("{}iByAm",s), &format!("{}iBiH",s)], [&format!("{}yE",s), &format!("{}iByAm",s), &format!("{}iByaH",s)], [&format!("{}yAH",s), &format!("{}iByAm",s), &format!("{}iByaH",s)], [&format!("{}yAH",s), &format!("{}yoH",s), &format!("{}InAm",s)], [&format!("{}yAm",s), &format!("{}yoH",s), &format!("{}izu",s)], [&format!("{}e",s), &format!("{}I",s), &format!("{}ayaH",s)]); }
    } else if base.ends_with("at") {
        let s2 = &base[..base.len()-2];
        if gender == "masculine" { res = build_grid([&format!("{}an",s2), &format!("{}antO",s2), &format!("{}antaH",s2)], [&format!("{}antam",s2), &format!("{}antO",s2), &format!("{}ataH",s2)], [&format!("{}atA",s2), &format!("{}adByAm",s2), &format!("{}adBiH",s2)], [&format!("{}ate",s2), &format!("{}adByAm",s2), &format!("{}adByaH",s2)], [&format!("{}ataH",s2), &format!("{}adByAm",s2), &format!("{}adByaH",s2)], [&format!("{}ataH",s2), &format!("{}atoH",s2), &format!("{}atAm",s2)], [&format!("{}ati",s2), &format!("{}atoH",s2), &format!("{}atsu",s2)], [&format!("{}an",s2), &format!("{}antO",s2), &format!("{}antaH",s2)]); }
    } else if base.ends_with("an") {
        let s2 = &base[..base.len()-2];
        if gender == "masculine" { res = build_grid([&format!("{}A",s2), &format!("{}AnO",s2), &format!("{}AnaH",s2)], [&format!("{}Anam",s2), &format!("{}AnO",s2), &format!("{}naH",s2)], [&format!("{}nA",s2), &format!("{}aByAm",s2), &format!("{}aBiH",s2)], [&format!("{}ne",s2), &format!("{}aByAm",s2), &format!("{}aByaH",s2)], [&format!("{}naH",s2), &format!("{}aByAm",s2), &format!("{}aByaH",s2)], [&format!("{}naH",s2), &format!("{}noH",s2), &format!("{}nAm",s2)], [&format!("{}ni",s2), &format!("{}noH",s2), &format!("{}asu",s2)], [&format!("{}an",s2), &format!("{}AnO",s2), &format!("{}AnaH",s2)]); }
    }
    
    if res.is_empty() {
        return Err(format!("Declension logic for '{}' in gender '{}' is not fully implemented in Rust yet.", base, gender));
    }
    
    for forms in res.values_mut() {
        for word in forms.iter_mut() {
            *word = apply_natva(word);
        }
    }
    
    Ok(res)
}
