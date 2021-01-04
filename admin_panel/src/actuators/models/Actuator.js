export default class Actuator {
    constructor(ip, mappingList) {
        this.ip = ip;
        this.mappingList = mappingList;
    }

    static isActuatorData(obj) {
        return obj.hasOwnProperty("ip");
    }

    static init(ip, mappingList, situationList) {
        console.log(situationList);
        console.log(mappingList);
        mappingList.forEach(mapping => {
            const targetSituationIndex = situationList.findIndex(situation => situation.id === mapping.situation);
            if (targetSituationIndex !== -1) {
                mapping.situation = situationList[targetSituationIndex].name;
            }
        });
        return new Actuator(ip, mappingList);
    }

    static mappingWithName(mapping, situationList) {
        const targetSituationIndex = situationList.findIndex(situation => situation.id === mapping.situation);
        mapping.situation = situationList[targetSituationIndex].name;
        return mapping;
    }
}
