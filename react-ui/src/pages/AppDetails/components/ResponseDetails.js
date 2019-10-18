import React from 'react';

const ResponseDetails = (response) => {
    return (
        <div>
            <div>{response.active}</div>
            <div>{response.body}</div>
            <div>{response.status_code}</div>
            <div>{response.header && response.header.map(({ name, value }) => (
                <>
                    Header:
                    <span>{name}</span>
                    <span>{value}</span>
                </>
            ))}</div>
        </div>
    )
};

export default ResponseDetails;
