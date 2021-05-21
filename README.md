# [Sandeadry (san-dead-ry)](https://choicez.alphaman.me/)
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

## Example file

```json
[
    {
        "Question": "1. The self organizing list improves the efficiency of _______",
        "Options": [
            "binary search",
            "jump search",
            "sublist search",
            "linear search"
        ],
        "Answer": "d",
        "Explanation": "Linear search in a linked list has time complexity O(n). To improve the efficiency of the linear search the self organizing list is used. A self-organizing list improves the efficiency of linear search by moving more frequently accessed elements towards the head of the list."
    },
    {
        "Question": "2. Which of the following is true about the Move-To-Front Method for rearranging nodes?",
        "Options": [
            "node with highest access count is moved to head of the list",
            "requires extra storage",
            "may over-reward infrequently accessed nodes",
            "requires a counter for each node"
        ],
        "Answer": "c",
        "Explanation": "In Move-To-front Method the element which is searched is moved to the head of the list. And if a node is searched even once, it is moved to the head of the list and given maximum priority even if it is not going to be accessed frequently in the future. Such a situation is referred to as over-rewarding."
    },
    {
        "Question": "3. What technique is used in Transpose method?",
        "Options": [
            "searched node is swapped with its predecessor",
            "node with highest access count is moved to head of the list",
            "searched node is swapped with the head of list",
            "searched nodes are rearranged based on their proximity to the head node"
        ],
        "Answer": "a",
        "Explanation": "In Transpose method, if any node is searched, it is swapped with the node in front unless it is the head of the list. So, in Transpose method searched node is swapped with its predecessor."
    }
]
```

## Running

Go to [releases page](https://github.com/smartclash/Sandeadry/releases) and download the executable file for your operating system.
    > Note: MacOS build is not available. I am broke, and I don't own a Mac


### To scrape Subjects, topics and it's MCQs

Rename the file into `sandeadry` and run the following. You can replace the link with a link to any other degree
```bash
./sandeadry scrape https://www.sanfoundry.com/computer-science-questions-answers/
```

### To save the data scrapped

All the data scrapped will be cleaned and saved into an `sqlite` database for further querying and/or manipulation.

```bash
./sandeadry save
```

### Index MCQs to Meilisearch

[Meilisearch](https://www.meilisearch.com/) is a search engine that provides results as you type. Make sure to run `save` command first before trying to index. 

1. Copy `.env.example` into `.env`. And enter the correct meilisearch instance credentials
    ```bash
    cp .env.example .env
    ```

2. Run the command and wait ;)
    ```bash
    ./sandeadry index
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
