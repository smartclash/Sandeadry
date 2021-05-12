# Sandeadry (san-dead-ry)
***A sanfoundry scrapper that does the heavylifting***

Just pass the link to Degree you want to scrape. Program will automatically create folders with subject names and save all the MCQs of that particular subject as a `json` file. Each `json` file belongs to a particular topic.

> **Example degree link: https://www.sanfoundry.com/computer-science-questions-answers/**

After you run the program, it will save all the MCQs in following format

```bash
.
└── datas
  └── degree_name
    ├── subject_name_1
    │ ├── topic_name_1.json
    │ └── topic_name_2.json
    └── subject_name_2
      ├── topic_name_1.json
      └── topic_name_2.json
```

## Requirements

This is a nodejs project, you need the following

- Node [(download)](https://nodejs.org/en/download/)
- Yarn [(download)](https://yarnpkg.com/getting-started/install)

## Installation

Everything is done via a terminal. So make sure to open yours.

1. Clone this project
```bash
git clone https://github.com/smartclash/Sandeadry.git
```

2. Get inside and install dependencies
```bash
yarn install
```

That's it. You've done the installation

## Running

Open the file `Index.ts` inside `Sandeadry/src` in your favorite text editor and replace `degreeName` with the name of the degree you're trying to scrape the MCQs off and replace `degreeLink` with the link to degree index page

Again in terminal, run. Make sure you're inside the project directory
```bash
yarn ts-build
```

Now to start scrapping, just run
```bash
node .
```

## Disclaimer

It's a fun project. I did it in my free time. Wanted to test out my knowledge and here we are.

This project is kinda buggy. Does the work 90% of the time but messes up the questions sometimes. If you can find the issue, fixes and PRs are highly appreciated.

## To Sanfoundry Folks

Just let me know and I'll remove this code. Don't DMCA lol, I am done with DMCA requests and can't handle anymore.
