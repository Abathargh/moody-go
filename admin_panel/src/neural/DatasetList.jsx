import React, {Component, useEffect, useState} from "react";
import socketIOClient from "socket.io-client";

const ENDPOINT = "http://moodybase:7000";


export default function DatasetList() {
    const [response, setResponse] = useState("");

    useEffect(() => {
        const socket = socketIOClient(ENDPOINT);
        socket.on("data", data => {
            setResponse(data);
        });
    }, []);

    return <p>{response}</p>;

}