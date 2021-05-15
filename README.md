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
    │  ├── topic_name_1.json
    │  └── topic_name_2.json
    └── subject_name_2
      ├── topic_name_1.json
      └── topic_name_2.json
```

## Running

1. Go to [releases page](https://github.com/smartclash/Sandeadry/releases) and download the executable file for your operating system.
    > Note: MacOS build is not available. I am broke, and I don't own a Mac

2. Rename the file into `sandeadry` and run the following. You can replace the link with a link to any other degree
    ```bash
    ./sandeadry -l https://www.sanfoundry.com/computer-science-questions-answers/
    ```

## Contributing

Project is based on [Golang](https://golang.org/). Clone the repo, and you can start working on it. If you have any problem, make an issue, and I'll get back to you

> Project was first built on [NodeJS (TypeScript)](https://github.com/smartclash/Sandeadry/tree/nodejs). It since has been re-written into golang.

## Disclaimer

It's a fun project. I did it in my free time. Wanted to test out my knowledge and here we are.

This project will never get a stable release and will always be a beta software. If you find any issue, please file that in [issue page](https://github.com/smartclash/Sandeadry/issues).
I'll appreciate it more if you can fix bugs and send PRs on my way. Do checkout [Contributing](#Contributing)

## To Sanfoundry Folks

Just let me know, and I'll remove this code. Don't DMCA lol, I am done with DMCA requests and can't handle anymore.
