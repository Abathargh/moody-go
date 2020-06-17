import React, {Component} from "react";
import  {Creator} from "./Activators";
import "./Creators.css"

export default class CreatorBox extends Component {
    constructor(props) {
        super(props);
        this.state = {
            newName: ""
        }

        this.handleChange = this.handleChange.bind(this);
        this.forwardData = this.forwardData.bind(this);
    }

    handleChange(evt) {
        this.setState({newName: evt.target.value});
    }

    forwardData() {
        const newName = this.state.newName;
        this.props.handleCreate(newName);
    }

    render() {
        return(
            <div className="creator">
                <table className="creatorTable">
                    <tbody>
                        <tr>
                            <td><h2>New {this.props.name}</h2></td>
                        </tr>
                        <tr>
                            <td><input type="text" className="nameInput" onChange={this.handleChange}/></td>
                        </tr>
                        <tr>
                            <td><Creator handleCreate={() => this.forwardData()}/></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        );
    }
}
