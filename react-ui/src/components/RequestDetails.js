import React from 'react';
import ResponseDetails from './ResponseDetails';
import BodyDetails from "./BodyDetails";

const RequestDetails = ({ app }) => {
    console.log('app.requests', app.requests);
    return (
        <div style={{border: '1px solid red', padding: '1em'}}>
            <h2>Request</h2>
            {app && app.requests && app.requests.map((request) => (
                <div key={`request-${request.id}`}>
                    <div><b>ID: </b>{request.id}</div>
                    <div><b>URL: </b>{request.url}</div>
                    <div><b>Method: </b>{request.method}</div>
                    <dl><b>Headers: </b>{request.header && request.header.map(({ name, value }) => (
                        <React.Fragment key={`request-${request.id}-header-${name}`}>
                            <dt><b>{name}</b></dt>
                            <dd>{value}</dd>
                        </React.Fragment>
                    ))}</dl>
                    <BodyDetails base64Body={request.body}/>
                    <div>{request.responses && request.responses.map(response => (
                        <ResponseDetails key={`request-${request.id}-response-${response.id}`} response={response}/>
                    ))}</div>
                </div>
            ))}
        </div>
    )
};

export default RequestDetails;
