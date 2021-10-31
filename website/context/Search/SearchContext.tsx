import { createContext, Dispatch, useContext, useReducer } from "react";
import searchReducer, { SearchAction, SearchContextState } from "./reducer";

interface SearchContextType {
  dispatchSearch: Dispatch<SearchAction>;
  searchState: SearchContextState;
}

const defaultState: SearchContextState = { search: null };

const SearchContext = createContext<SearchContextType | undefined>(undefined);

const SearchProvider: React.FC = ({ children }) => {
  const [searchState, dispatchSearch] = useReducer(searchReducer, defaultState);

  return (
    <SearchContext.Provider value={{ searchState, dispatchSearch }}>
      {children}
    </SearchContext.Provider>
  );
};

const useSearch = (): SearchContextType => {
  const context = useContext(SearchContext);
  if (context === undefined) {
    throw new Error("useSearch must be used within a SearchProvider");
  }
  return context;
};

export { SearchProvider, useSearch };
