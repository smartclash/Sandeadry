import { parametreize } from './Helpers'
import * as writeJsonFile from 'write-json-file'
import MCQParser, { MCQ } from './Parser/MCQParser'
import SubjectParser, { Topic } from './Parser/SubjectParser'
import DegreeParser, { Subject } from './Parser/DegreeParser'

const degreeName = 'Computer Science'
const degreeLink = 'https://www.sanfoundry.com/computer-science-questions-answers/'

const writeJson = (degreeName: string, subject: Subject, topic: Topic, mcqs: MCQ[]) => {
    const p = parametreize
    const fileName = p(topic.name) + '.json'
    const filePath = `datas/${p(degreeName)}/${p(subject.name)}/${fileName}`

    return writeJsonFile.sync(filePath, mcqs, { indent: 2 })
}

const handler = async () => {
    const degreeParser = await DegreeParser(degreeName, degreeLink)
    degreeParser.subjects.forEach(async (subjects) => {
        const subjectName = subjects.name
        const subjectParser = await SubjectParser(subjectName, subjects.link)

        subjectParser.topics.forEach(async (topics) => {
            const topicName = topics.name
            const mcqParser = await MCQParser(topicName, topics.link)

            writeJson(degreeName, subjects, topics, mcqParser.mcqs)
        })
    })
}

handler()
