# 📊 GoLog Analyzer — *loganalyzer*

CHEIO Thomas - Mario Aboujamra

Un outil **CLI en Go** pour analyser des fichiers de logs en **parallèle** avec :  
✔️ Concurrence (goroutines + WaitGroup)  
✔️ Erreurs personnalisées (`FileAccessError`, `LogParseError`)  
✔️ Import JSON (config) & Export JSON (résultats)  
✔️ Interface CLI avec **Cobra** (`analyze`, `add-log`)  

---

## 🚀 Installation

```bash
# 1) Clone du projet
git clone https://github.com/tcheio/GOLogAnalyser.git
cd GOLogAnalyser

# 2) Initialiser les dépendances
go mod tidy

# 3) Vérifier que Cobra est présent
go get github.com/spf13/cobra@v1.8.0
```

---

## 📁 Préparer les fichiers de test

Dans le dossier `test_logs/`, ajoute :  
- **access.log** → log type *nginx/access*  
- **errors.log** → log applicatif avec des erreurs/warnings/fatals  
- **empty.log** → fichier vide (test d’erreur)  
- **corrupted.log** → fichier volontairement corrompu (test parsing)  

Crée un fichier `config.json` :

```json
[
  { "id": "access-log",    "path": "./test_logs/access.log",    "type": "nginx-access" },
  { "id": "errors-log",    "path": "./test_logs/errors.log",    "type": "custom-app" },
  { "id": "empty-log",     "path": "./test_logs/empty.log",     "type": "generic" },
  { "id": "corrupted-log", "path": "./test_logs/corrupted.log", "type": "generic" }
]
```

---

## ▶️ Utilisation

### 🔎 Analyse simple (résultats console)
```bash
go run . analyze -c ./config.json
```

### 📤 Analyse + Export JSON (dans `./reports/`)
```bash
go run . analyze -c ./config.json -o ./reports/report.json
```
➡️ Génère `./reports/YYMMDD_report.json` (préfixé par la date).

### 🎯 Filtrer par statut
```bash
# Afficher uniquement les résultats OK
go run . analyze -c ./config.json --status OK

# Afficher uniquement les résultats FAILED
go run . analyze -c ./config.json --status FAILED
```

### ➕ Ajouter un log à la configuration
```bash
go run . add-log \\
  --id newlog \\
  --path ./test_logs/new.log \\
  --type generic \\
  --file ./config.json
```

---

## ✅ Exemple attendu (console)

```
--- Résumé de l'analyse ---
[access-log] ./test_logs/access.log — OK — access.log: 120 lignes (2xx=115, 3xx=2, 4xx=3, 5xx=0)
[errors-log] ./test_logs/errors.log — OK — errors.log: 50 lignes (error=20, warn=25, fatal=5)
[empty-log] ./test_logs/empty.log — FAILED — Fichier vide.
[corrupted-log] ./test_logs/corrupted.log — FAILED — Erreur de parsing.
  ↳ erreur: erreur de parsing dans './test_logs/corrupted.log': contenu corrompu détecté
```

---

## 🛠 Développement

### Compiler le binaire
```bash
go build -o loganalyzer
```

### Lancer après compilation
```bash
./loganalyzer analyze -c ./config.json
```

---
