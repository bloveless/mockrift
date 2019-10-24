import React from 'react';
import AceEditor from 'react-ace';

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-github";

const BodyDetails = ({ id, contentType, base64Body }) => {
    const body = parseBase64Body(contentType, base64Body);
    return (
        <AceEditor
            mode="json"
            theme="github"
            name={`body-details-ace-editor-${id}`}
            value={body}
            editorProps={{ $blockScrolling: true }}
        />
    )
};

const parseBase64Body = (contentType, body) => {
    let parsedBody = atob(body);
    if (contentType === 'application/json') {
        parsedBody = JSON.stringify(JSON.parse(parsedBody), null, 2);
        parsedBody = parsedBody.replace(/\\n/g, "\n")
    }

    return parsedBody;
};

export default BodyDetails;