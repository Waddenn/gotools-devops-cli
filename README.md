# GoTools

[![CI](https://github.com/Waddenn/gotools-devops-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/Waddenn/gotools-devops-cli/actions/workflows/ci.yml)

Outil en ligne de commande en Go pour le projet M1 DevOps.

Il regroupe des fonctionnalités de manipulation de fichiers texte, récupération de contenu Wikipédia, gestion de processus, sécurité (verrouillage / permissions) et vérifications d'infra.

## Niveau vise

18/20

## Lancer le projet

```bash
go mod tidy
go build -o gotools .
./gotools
```

Optionnel : charger un fichier de config précis.

```bash
./gotools --config config.json
# ou
./gotools --config config.txt
```

## Menus disponibles

### Fonctionnalites implementees

- `A` : analyser un fichier texte (infos, stats, filtrage, head/tail)
- `B` : analyser plusieurs fichiers `.txt` d'un dossier (rapport, index, fusion)
- `C` : récupérer un ou plusieurs articles Wikipedia et produire des stats simples
- `D` : lister / rechercher / arreter un processus (avec confirmation)
- `E` : verrouiller un fichier (lockfile), changer les permissions, vérifier les droits
- `F` : afficher les conteneurs Docker actifs + stats d'un conteneur
- `G` : vérifier l'espace disque restant
- `H` : scanner plusieurs fichiers en parallèle (goroutines + `WaitGroup`)

## Compatibilite OS

Le projet est teste en CI sur `ubuntu-latest`, `macos-latest` et `windows-latest`.

- `A` / `B` (FileOps) : compatible Linux/macOS/Windows
- `C` (WebOps) : compatible Linux/macOS/Windows (acces reseau requis)
- `D` (ProcOps) : commandes adaptees selon l'OS
  - Windows : `tasklist`, `taskkill`
  - macOS/Linux : `ps`, `kill`
- `E` (SecureOps) : lockfile portable + changement de permissions
- `F` (ContainerOps) : compatible Linux/macOS/Windows si Docker CLI est installe et actif
- `G` (Etat disque) : compatible Linux/macOS/Windows (`df` sur Unix, `wmic`/PowerShell sur Windows)
- `H` (Scan parallele) : compatible Linux/macOS/Windows

En cas d'outil systeme manquant (ex: Docker non installe), le programme retourne une erreur claire.

## Structure du projet

```text
main.go                 menu principal
config/config.go        chargement config (txt/json)
fileops/analysis.go     analyse d'un fichier
fileops/multi.go        opérations sur plusieurs fichiers
webops/wiki.go          récupération / analyse Wikipedia
procops/process.go      gestion des processus
secureops/secure.go     lockfile + permissions
infraops/container.go   infos Docker
infraops/health.go      vérification espace disque
audit/audit.go          journalisation des actions sensibles
```

## Fichiers utiles

- `config.json` / `config.txt` : configuration
- `data/` : exemples de fichiers d'entrée
- `out/` : fichiers générés (rapports, filtres, logs)

## Remarques

- Le menu est volontairement simple et lisible, sans framework CLI.
- Certaines fonctions dépendent de l'environnement (`docker`, `ps`, `df`, etc.).
- Les actions sensibles (arret de processus, lock, chmod) sont tracees dans `out/audit.log`.
- Une CI GitHub Actions a ete ajoutee (verification format/build/vet + smoke test CLI + controle des livrables) avec execution sur tags de release (`v*`, `release-*`) et declenchement manuel.

## Description du travail effectue

- Decoupage du projet en packages par domaine (`fileops`, `webops`, `procops`, `secureops`, `infraops`).
- Menu CLI interactif en boucle avec configuration `config.txt` et `config.json` (`--config`).
- Gestion des erreurs renforcee sur les traitements fichiers et les operations Docker/OS.
- Ajout des operations sensibles avec confirmation (`kill`, lock/unlock) + journalisation dans `out/audit.log`.
- Ajout de traitements paralleles avec goroutines + `sync.WaitGroup` + channels (scan fichiers et Wikipedia multi-articles).
