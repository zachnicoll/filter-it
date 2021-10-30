import { Progress } from "api/types";
import { toastUpdate } from "common/toast";
import {
  Dispatch,
  createContext,
  useEffect,
  useReducer,
  useRef,
  useContext,
} from "react";
import progressReducer, {
  ProgressAction,
  ProgressContextState,
} from "./reducer";

interface ProgressContextType {
  dispatchProgress: Dispatch<ProgressAction>;
  progressState: ProgressContextState;
}

const defaultState: ProgressContextState = {};

const ProgressContext = createContext<ProgressContextType | undefined>(
  undefined
);

const ProgressProvider: React.FC = ({ children }) => {
  const [progressState, dispatchProgress] = useReducer(
    progressReducer,
    defaultState
  );

  const intervalRef = useRef<any>(undefined);

  useEffect(() => {
    if (intervalRef.current && progressState.status == Progress.PROCESSING) {
      intervalRef.current = setInterval(() => {
        if (progressState.notifyReference) {
          toastUpdate(progressState.notifyReference as string);
        }
      }, 15000);
    } else if (
      intervalRef.current &&
      progressState.status !== Progress.PROCESSING
    ) {
      clearInterval(intervalRef.current);
    }
  }, [progressState]);

  return (
    <ProgressContext.Provider value={{ progressState, dispatchProgress }}>
      {children}
    </ProgressContext.Provider>
  );
};

const useProgress = (): ProgressContextType => {
  const context = useContext(ProgressContext);
  if (context === undefined) {
    throw new Error("useProgress must be used within a ProgressProvider");
  }
  return context;
};

export { ProgressProvider, useProgress };
