# StratMind

Open-source AI Tactical Coach for CS2 demos

StratMind/
├── README.md # ✨ Manifest + scop + cum contribui
├── LICENSE # 🧾 MIT License (sau altă licență)
├── .gitignore # Ignoră fișiere temporare
├── go.mod / go.sum # Dependințe Go
│
├── core/ # 🧠 Tot ce e open-source și face parsing logic
│ ├── parser/ # Folosește demoinfocs-golang (citirea demo-urilor)
│ ├── rules/ # Reguli tactice (ex: solo B, greșeli rotație)
│ │ ├── mirage.go
│ │ ├── inferno.go
│ │ └── nuke.go
│ └── utils/ # Funcții comune (distanță, unghiuri, conversii)
│
├── prompts/ # 📝 Prompturi GPT generate de reguli
│ ├── mirage/ # Organizate pe hartă
│ │ └── b_entry_solo.md
│ └── common/ # Prompturi generale (peek greșit, rotation etc.)
│
├── cmd/ # ⚙️ Tool-uri CLI pentru testat (dev only)
│ └── stratmind.go
│
├── webapp/ # 🌐 Interfață UI (open-core)
│ ├── frontend/ # React / Tailwind (sau alt framework web)
│ └── backend/ # API GPT, auth, planuri
│
├── docs/ # 📚 Documentație: dev guide, cum rulezi, cum contribui
│ ├── architecture.md
│ ├── contribution-guide.md
│ └── how-it-works.md
│
├── roadmap.md # 📈 Planul de evoluție al proiectului
├── SUPPORT.md # 💬 Unde ceri ajutor / cum raportezi buguri
├── CONTRIBUTING.md # 🫱 Cum contribui cu cod, reguli, prompturi
├── .env.example # ⚠️ Config demo (.env fără chei reale)
└── images/ # 🖼️ Logo, screenshots, badges pentru README

## 🧾 License

This project is licensed under the [MIT License](LICENSE).
