export const parametreize = (string: string) =>
    string.replace(/\s/gu, '_').replace('/', '_OR_').toLowerCase()
