import { Filter } from "api/types";
import produce from "immer";

export type SearchAction = { type: "SEARCH"; payload: Filter | null };

export interface SearchContextState {
  search: Filter | null;
}

const searchReducer = produce(
  (state: SearchContextState, action: SearchAction) => {
    switch (action.type) {
      case "SEARCH":
        state.search = action.payload;
        break;
    }
  }
);

export default searchReducer;
