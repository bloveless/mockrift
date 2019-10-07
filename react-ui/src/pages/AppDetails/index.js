import React, { useState } from 'react';

const AppDetails = ({ match }) => {
    return <div>App Details for {match.params.slug}</div>;
};

export default AppDetails;
