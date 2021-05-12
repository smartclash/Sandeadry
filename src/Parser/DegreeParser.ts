import { JSDOM } from 'jsdom'

interface Subject {
    name: string,
    link: string
}

interface SubjectParserResult {
    degree: string,
    total: number,
    subjects: Subject[]
}

interface SubjectParser {
    (degree: string, url: string): Promise<SubjectParserResult>
}

const DegreeParser: SubjectParser = async (degree, url) => {
    const $ = await JSDOM.fromURL(url)
    const doc = $.window.document

    let subjectsRaw: SubjectParserResult = { degree, total: 0, subjects: [] }

    doc.querySelectorAll('div.entry-content > table > tbody > tr > td > li > a').forEach(element => {
        const name = element.textContent
        if (!name.toLowerCase().includes('tests'))
            return subjectsRaw.subjects.push({
                name: element.textContent,
                //@ts-ignore
                link: element.href
            })
    })

    subjectsRaw.total = subjectsRaw.subjects.length
    return subjectsRaw
}

export default DegreeParser
