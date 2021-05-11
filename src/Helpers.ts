export const parametreize = (string: string) => {
    return string.replace(/\s/gu, '_').toLowerCase()
};
