import React from "react"

import Back from "../../../components/Card/Back"
import Face from "../../../components/Card/Face"
import Suite from "../../../components/Card/Suite"

export default function Deck({trump_card, amount}) {
    return <div className="deck">
        {amount >= 2 && <Back className="top-card" />}
        {amount >= 1 && <Face className="main-suite" card={trump_card} />}
        {amount || <Suite suite={trump_card.suite} />}
    </div>
}