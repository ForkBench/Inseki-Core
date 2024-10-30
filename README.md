<a name="readme-top"></a>

<div align="center">
  <img src="doc/Inseki.png" alt="logo" width="140" height="auto" />
  <br/>

<h3><b>Inseki - Project Discovery Tool</b></h3>

</div>

# 📗 Table of Contents

- [📗 Table of Contents](#-table-of-contents)
- [📖 Inseki ](#-inseki-)
  - [🛠 Built With ](#-built-with-)
    - [Tech Stack ](#tech-stack-)
    - [Key Features ](#key-features-)
  - [💻 Getting Started ](#-getting-started-)
    - [Prerequisites](#prerequisites)
    - [Setup](#setup)
    - [Build](#build)
    - [Development](#development)
    - [Example](#example)
  - [🔭 Future and Current Features ](#-future-features-)
  - [📝 License ](#-license-)

# 📖 Inseki <a name="about-project"></a>

**Inseki** is an Open-Source tool designed for discovering and analyzing project structures within a disk. It scans directories and represents the structure of each project in JSON format.

👷‍ ~~It is currently under development and is not yet ready for production use.~~

This project is now working ! You can use it to scan your disk and see the results.

## 🛠 Built With <a name="built-with"></a>

### Tech Stack <a name="tech-stack"></a>

The project is built using the following technology:

<details>
  <summary>Back-End</summary>
  <ul>
    <li><a href="https://go.dev">Go</a></li>
  </ul>
</details>

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Key Features <a name="key-features"></a>

- 🚀 Scan directories to discover project structures
- 🗂 Represent project structures in JSON format

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 💻 Getting Started <a name="getting-started"></a>

To get a local copy up and running, follow these steps.

### Prerequisites

In order to run this project you need:

- [Go](https://golang.org/dl/)

### Setup

Clone this repository to your desired folder:

```
  cd my-folder
  git clone git@github.com:ForkBench/Inseki-Core.git
```

### Build

Build the project using the following command:

```
  go build
```

To run it, execute the following command:

```
  ./inseki-core
```

You'll have to put structures into `~/.inseki` directory to see the results.

### Development

To run the project, execute the following command:

```
  go run .
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Example

A config :

File located at : `~/.inseki/.insekiignore` to avoid some folders.

```gitignore
node_modules
.bun
.local
.git
.idea
.vscode
pkg
libraries
```

File located at : `~/.inseki/structures/C-programming/projects.json` to define some C projects.

```json
{
    "name": "*",
    "isDirectory": true,
    "children": [
        {
            "name": "src",
            "isDirectory": true,
            "children": [
                {
                    "name": "*.c",
                    "isDirectory": false
                }
            ]
        },
        {
            "name": "lib",
            "isDirectory": true,
            "children": [
                {
                    "name": "*.h",
                    "isDirectory": false
                }
            ]
        }
    ]
}
```

Example of output : 

```bash
$ go run .
Number of structures analysed: 3
Number of files analysed: 13739
Filepath: .../courses/S1/C/TP-Temp/TP1 - Outils/Part2/teZZt.h, Structure: lab.json
Filepath: .../revisions-c/Exercice/Tri insertion/main.h, Structure: lab.json
Filepath: .../revisions-c/Exercice/Tri insertion/main.c, Structure: lab.json
Filepath: .../courses/S1/C/TP-Temp/TP1 - Outils/Part2/exemple.c, Structure: lab.json
...
```


<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🔭 Future and Current Features <a name="future-features"></a>

- [x] Add a structure library for common project types
- [x] Add go routines for faster scanning
- [x] Add a way to ignore some folders
- [ ] Add advanced filtering options for project discovery
- [ ] Implement a graphical interface for easier interaction
- [ ] Remove multiple scan of files (if a project is composed with n files, it will scan and validate n times)
- [ ] Test more the project

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 📝 License <a name="license"></a>

This project is [GNU](LICENSE) licensed.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
