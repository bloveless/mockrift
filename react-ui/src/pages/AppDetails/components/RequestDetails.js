import React from 'react';
import ResponseDetails from './ResponseDetails';

const RequestDetails = ({ requests }) => {
    console.log(requests);
    return (
        <div>
            {requests && requests.map((request) => (
                <div>
                    <div>{request.url}</div>
                    <div>{request.method}</div>
                    <div>{request.header && request.header.map(({ name, value}) => (
                        <>
                            Header:
                            <span>{name}:</span>
                            <span>{value}</span>
                        </>
                    ))}</div>
                    <div>{request.body}</div>
                    <div>{request.responses && request.responses.map(response => (
                        <ResponseDetails response={response}/>
                    ))}</div>
                </div>
            ))}
        </div>
    )
};

export default RequestDetails;
