import React, {Component} from "react";
import {ErrorMessage} from "./Error";
import Situation from "./Situation";
import "./List.css";

export default class SituationList extends Component {
    constructor(props) {
        super(props);
        this.removeSituation = this.removeSituation.bind(this);
    }

    removeSituation(id) {
        this.props.handleSituationRemoval(id);
    }

    render() {
        if(this.props.situationList.length === 0) {
            return(
                <div className="list">
                    <h2>Situations</h2>
                    <ErrorMessage name="No situations!" />
                </div>
            );
        }else{
            return(
                <div className="situationList">
                    <h2>Situations</h2>
                    <table className="listTable">
                        <thead>
                        <tr>
                            <th>Name</th>
                            <th>Activate</th>
                        </tr>
                        </thead>
                        <tbody>
                        { this.props.situationList.map(situation => (
                            <Situation key={situation.id}
                                       id={situation.id}
                                       name={situation.name}
                                       handleRemove={this.removeSituation}/>
                        ))}
                        </tbody>
                    </table>
                </div>
            );
        }
    }
}