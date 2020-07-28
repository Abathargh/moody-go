export default class Actuator {
    static init(ip, mappingList, situationList) {
        mappingList.forEach(mapping => {
            const targetSituationIndex = situationList.findIndex(situation => situation.id === mapping.situation);
            mapping.situation = situationList[targetSituationIndex].name;
        });
        return new Actuator(ip, mappingList);
    }

    static mappingWithName(mapping, situationList) {
        const targetSituationIndex = situationList.findIndex(situation => situation.id === mapping.situation);
        mapping.situation = situationList[targetSituationIndex].name;
        return mapping;
    }

    constructor(ip, mappingList) {
        this.ip = ip;
        this.mappingList = mappingList;
    }
}
