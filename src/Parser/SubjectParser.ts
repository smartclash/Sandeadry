import { JSDOM } from 'jsdom'

interface Topic {
    name: string,
    link: string
}

interface SubjectParserResult {
    subject: string,
    topics: Topic[]
}

interface SubjectParser {
    (subject: string, link: string): Promise<SubjectParserResult>
}

const SubjectParser: SubjectParser = async (subject, link) => {
    const $ = await JSDOM.fromURL(link)
    const document = $.window.document

    let topicsObject: SubjectParserResult = { subject, topics: []}

    document.querySelectorAll('div.sf-section > table > tbody > tr').forEach(tr => {
        tr.childNodes.forEach(td => {
            td.childNodes.forEach(li => {
                li.childNodes.forEach(a => {
                    topicsObject.topics.push({
                        name: a.textContent,
                        //@ts-ignore
                        link: a.href
                    })
                })
            })
        })
    })

    return topicsObject
}

export default SubjectParser
