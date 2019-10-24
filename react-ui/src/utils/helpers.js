const getContentTypeValue = (header) => {
    const contentTypeHeader = header && header.find((header) => header.name === 'Content-Type');
    return contentTypeHeader && contentTypeHeader.value.length > 0 && contentTypeHeader.value[0].split(';')[0];
};

export {
    getContentTypeValue,
}
