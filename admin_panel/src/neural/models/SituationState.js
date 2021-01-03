export default class SituationState {
    constructor(isSet, id, name) {
        this.isSet = isSet;
        this.id = id;
        this.name = name;
    }

    sameId(otherId) {
        return parseInt(otherId, 10) === parseInt(this.id, 10)
    }

    static FromAPI(apiResp) {
        if (!apiResp.isSet) {
            return new SituationState(false, null, "None");
        }
        let situation = apiResp.situation;
        return new SituationState(true, situation.id, situation.name);
    }
}