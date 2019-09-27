import React from "react";
import ReactDOM from "react-dom";
import { getSettings } from './settings';

const App = () => {
    const settings = getSettings(process.env.MODE);

    return <div>Hello React,Webpack 4 & Babel 7! {settings.API_URL}</div>;
};

ReactDOM.render(<App/>, document.querySelector("#root"));
