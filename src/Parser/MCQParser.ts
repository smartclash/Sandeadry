import { JSDOM } from 'jsdom'

interface MCQ {
    question: string,
    options: string[],
    answer: string,
    explanation: string
}

interface MCQParserResult {
    topic: string,
    mcqs: MCQ[]
}

interface MCQParser {
    (topic: string, mcqLink: string): Promise<MCQParserResult>
}

let optionsObject = []
let answersObject = []
let rawQuestionsObject = []
let questionsObject = []

const optionsParser = (object: any[]): string[] => {
    let options = []

    object.forEach((value: string) => {
        const optionStarters = ['a)', 'b)' , 'c)', 'd)']
        optionStarters.forEach(starter => {
            if (!value.startsWith(starter))
                return

            options.push(value.replace(starter + ' ', ''))
        })
    })

    return options
}

const MCQParser: MCQParser = async (topic, mcqLink) => {
    const $ = await JSDOM.fromURL(mcqLink)
    const document = $.window.document

    document.querySelectorAll('div.entry-content > div.collapseomatic_content').forEach(element => {
        const ansObject = element.textContent.split('\n')
        const answer = ansObject[0].replace('Answer: ', '')
        const explanation = ansObject[1].replace('Explanation: ', '')

        answersObject.push({ answer, explanation })
    })

    const questionsDOM = document.querySelectorAll('div.entry-content, p')
    questionsDOM.forEach((element, key) => {
        if (key == 1 || key == 0)
            return

        let questionFound: boolean = false
        const qObject = element.textContent.split('\n')

        while (!questionFound) {
            if (qObject.includes('View Answer')) {
                rawQuestionsObject.push(qObject[0])
                questionFound = true
            }

            if (qObject.length == 1) {
                rawQuestionsObject.push(qObject[0])
                questionFound = true
            }
 
            questionFound = true
        }

        const parsedOptions = optionsParser(qObject)

        if (parsedOptions.length == 0)
            return
        
        optionsObject.push(parsedOptions)        
    })

    questionsObject = rawQuestionsObject
        .splice(0, rawQuestionsObject.length - 4)
        .filter((question: string) => !question.startsWith('a)'))

    let constructedMCQ: MCQParserResult = { topic, mcqs: [] }
    questionsObject.forEach((question: string, key: number) => {
        constructedMCQ.mcqs.push({
            question,
            options: optionsObject[key],
            answer: answersObject[key].answer,
            explanation: answersObject[key].explanation
        })
    })

    return constructedMCQ
}

export default MCQParser
