import React, { useState } from 'react';
import AceEditor from 'react-ace';

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-github";
import { UPDATE_REQUEST_BODY } from "../utils/queries";
import { useMutation } from "@apollo/react-hooks";

const BodyDetails = ({ id, contentType, base64Body }) => {
    const [body, updateBody] = useState(parseBase64Body(contentType, base64Body));
    const [editing, updateEditing] = useState(false);
    const [updateRequestBody] = useMutation(UPDATE_REQUEST_BODY);

    const saveBody = () => {
        console.log('New Body: ', btoa(body));
        updateRequestBody({
            variables: {
                id,
                body: btoa(body),
            },
        })
            .then(console.log)
            .catch((err) => console.log('Error', err));
    };

    if (!editing) {
        return (
            <div>
                <b>Body (
                    <button onClick={() => updateEditing(true)}>Toggle Edit</button>
                    )</b>
                <div>{body}</div>
            </div>
        );
    }

    return (
        <div>
            <b>Body ( <button onClick={() => updateEditing(false)}>Toggle Edit</button> )</b>
            <div>
                <AceEditor
                    mode="json"
                    theme="github"
                    name={`body-details-ace-editor-${id}`}
                    value={body}
                    editorProps={{ $blockScrolling: true }}
                    width="100%"
                    onChange={(newBody) => updateBody(newBody)}
                />
            </div>
            <div>
                <button onClick={saveBody}>Save Body</button>
            </div>
        </div>
    )
};

const parseBase64Body = (contentType, body) => {
    let parsedBody = atob(body);
    console.log(parsedBody);
    if (contentType === 'application/json') {
        parsedBody = JSON.stringify(JSON.parse(parsedBody), null, 2);
        parsedBody = parsedBody.replace(/\\n/g, "\n")
    }

    return parsedBody;
};

export default BodyDetails;