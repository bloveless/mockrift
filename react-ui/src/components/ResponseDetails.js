import React from 'react';
import BodyDetails from "./BodyDetails";
import { getContentTypeValue } from "../utils/helpers";

const ResponseDetails = ({ response }) => {
    return (
        <div style={{ border: '1px solid black' }}>
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
            <BodyDetails contentType={getContentTypeValue(response.header)} id={response.id} base64Body={response.body}/>
        </div>
    )
};

export default ResponseDetails;
