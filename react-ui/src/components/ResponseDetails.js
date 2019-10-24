import React from 'react';
import BodyDetails from "./BodyDetails";

const ResponseDetails = ({ response }) => {
    const contentTypeHeader = response.header && response.header.find((header) => header.name === 'Content-Type');
    const contentTypeValue = contentTypeHeader && contentTypeHeader.value.length > 0 && contentTypeHeader.value[0].split(';')[0];

    return (
        <div style={{ border: '1px solid black', padding: '1em' }}>
            <h2>Response</h2>
            <div><b>ID: </b>{response.id}</div>
            <div><b>Active: </b>{response.active}</div>
            <div><b>Status Code: </b>{response.status_code}</div>
            <div>{response.header && response.header.map(({ name, value }) => (
                <React.Fragment key={`response-${response.id}-header-${name}`}>
                    <dt><b>{name}</b></dt>
                    <dd>{value}</dd>
                </React.Fragment>
            ))}</div>
            <BodyDetails contentType={contentTypeValue} base64Body={response.body}/>
        </div>
    )
};

export default ResponseDetails;
