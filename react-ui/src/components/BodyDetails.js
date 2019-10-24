import React from 'react';

const BodyDetails = ({ base64Body }) => {
    return (
        <div><b>Body: </b>{atob(base64Body)}</div>
    )
};

export default BodyDetails;