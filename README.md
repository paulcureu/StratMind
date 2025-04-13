# 🧠 StratMind

> Primul AI Tactical Coach open-source pentru CS2 – începem cu Mirage.

StratMind analizează demo-urile `.dem` din meciurile CS2, detectează greșeli tactice și generează feedback inteligent cu ajutorul unui model de limbaj (ex: GPT).  
Scopul: să îți înțelegi mai bine greșelile și să progresezi ca jucător.  
Proiectul este complet open-source și construit împreună cu comunitatea.

---

## 🗺️ MVP Focus

🎯 Harta: **Mirage**  
📌 Primele reguli:

- entry solo fără flash
- lipsă sincronizare
- fake ineficient
- lipsă utilitare pe poziții cheie

---

## ⚙️ Cum funcționează

1. 📂 Tu încarci un `.dem` (demo de la meci)
2. 🧠 Parserul extrage tick-uri, poziții, utilitare, kills etc.
3. 📏 Se aplică reguli tactice (în Go)
4. 📝 Se generează un prompt contextual pentru GPT
5. 💬 Primești feedback detaliat, ca de la un coach adevărat

---

## 💻 Cum rulezi local

### 1. Clonează repo-ul

```bash
git clone https://github.com/tu/StratMind.git
cd StratMind
```
