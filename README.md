# TP : GoLog Analyzer - Analyse de Logs Distribuée

**golangWatcher** est un outil CLI en Go pour analyser des fichiers de logs. Il lit une liste de logs depuis un fichier JSON, les analyse en parallèle et génère un rapport au format JSON.

---

## Table des matières

- [Installation](#installation)
- [Structure du projet](#structure-du-projet)
- [Usage](#usage)
- [Exemple de JSON d'entrée](#exemple-de-json-dentrée)
- [Résultat attendu](#résultat-attendu)

---

## Installation

1.  Go installé dans votre pc :
```bash
go version
```

2.  Clonez le projet :
```bash
git clone https://github.com/BanggEddy/golangWatcher.git
```

3. Aller dans le répertoire :
```bash
cd golangWatcher
```

4. Installer les dépendances :
```bash
go mod tidy
```

5. Exécuter le projet :
```bash
go run main.go analyze --input config.json --output reports/report.json
```

## Structure du projet 

```bash
golangWatcher/
├─ cmd/
│  ├─ root.go         
│  ├─ analyze.go     
├─ internal/
│  ├─ analyzer/
│  │  ├─ analyze.go    
│  │  └─ errors.go    
│  ├─ config/
│  │  └─ config.go     
│  └─ reporter/
│     └─ report.go     
├─ test_logs/          
├─ reports/            
├─ config.json         
├─ go.mod
├─ go.sum
├─ main.go
└─ README.md
```

## Usage 

```bash
# Commande principale
go run main.go analyze --input <fichier_json_entree> --output <fichier_json_sortie>
```

Options :
| Option        | Description                                                      |
|---------------|------------------------------------------------------------------|
| `-i, --input` | Chemin vers le fichier JSON contenant la liste des logs.         |
| `-o, --output`| Chemin vers le fichier JSON pour exporter les résultats. (optionnel) |

## Exemple de JSON d'entrée (config.json) :
```bash
[
  {"id":"web-server-1", "path":"test_logs/access.log", "type":"info"},
  {"id":"app-backend-2", "path":"test_logs/errors.log", "type":"error"},
  {"id":"corrupted-log", "path":"test_logs/corrupted.log", "type":"error"},
  {"id":"db-server-3", "path":"test_logs/mysql_error.log", "type":"error"},
  {"id":"invalid-path", "path":"/non/existent/log.log", "type":"error"},
  {"id":"empty-log", "path":"test_logs/empty.log", "type":"info"}
]
```

## Résultat attendu

Lors de l'exécution :
```bash
go run . analyze -i config.json -o reports/report.json
```

On obtient :
```bash
❌ corrupted-log (test_logs/corrupted.log) : Erreur de parsing. - erreur de parsing sur test_logs/corrupted.log: ligne corrompue
❌ db-server-3 (test_logs/mysql_error.log) : Fichier introuvable. - fichier introuvable: test_logs/mysql_error.log (GetFileAttributesEx test_logs/mysql_error.log: The system cannot find the file specified.)
❌ invalid-path (/non/existent/log.log) : Fichier introuvable. - fichier introuvable: /non/existent/log.log (GetFileAttributesEx /non/existent/log.log: The system cannot find the path specified.)
✅ web-server-1 (test_logs/access.log) : Analyse terminée avec succès.
✅ app-backend-2 (test_logs/errors.log) : Analyse terminée avec succès.
✅ empty-log (test_logs/empty.log) : Analyse terminée avec succès.
✅ Résultats exportés vers reports/report.json
```
Le fichier JSON reports/report.json contiendra le résumé complet de l'analyse.


