import React from "react"
import {SearchResult} from "./SearchResult"
import "./SearchResultsList.css"

export const SearchResultsList = ({results}) => (
    <div className="results-list">
        <ol>
            {
                results.map((result, id) => <SearchResult result={result} key={id} />)
            }
        </ol>
    </div>
)
