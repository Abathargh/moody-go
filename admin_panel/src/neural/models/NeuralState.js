// An enum representing the three possible states of
// the dataset engine running on the gateway
const Mode = Object.freeze({
    "Stopped": "Stopped",
    "Collecting": "Collecting",
    "Predicting": "Predicting"
})

const stateMapping = [
    Mode.Stopped,
    Mode.Collecting,
    Mode.Predicting
]


export default class NeuralState {
    constructor(mode, dataset) {
        this.mode = mode;
        this.dataset = dataset;
    }

    static FromAPI(apiResp) {
        let mode = stateMapping[apiResp.mode];
        let dataset = apiResp.dataset === "" ? "None" : apiResp.dataset;
        return new NeuralState(mode, dataset);
    }

    toReqFormat() {
        let mode = stateMapping.indexOf(this.mode);
        let dataset = this.dataset === "None" ? "" : this.dataset;
        return { mode: mode, dataset: dataset };
    }

    equals(other) {
        return this.mode === other.mode && this.dataset === other.dataset;
    }
}

export { Mode };