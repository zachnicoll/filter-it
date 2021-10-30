import React, {
  useReducer,
  createContext,
  Dispatch,
  PropsWithChildren,
  ReactText,
  useEffect,
  useRef,
} from "react";
import { toast } from "react-toastify";

interface ProgressContextType {
  dispatch: Dispatch<ProgressAction>;
  state: ProgressState;
}

interface ProgressState {
  id?: string;
  status?: string;
  notifyReference?: ReactText;
  imageLocation?: string;
}

type ProgressAction =
  | { type: "UPDATE_IMAGE"; payload: string }
  | { type: "UPDATE_PROGRESS"; payload: string };

export const ProgressContext = createContext<ProgressContextType>(
  {} as ProgressContextType
);

const reducer = (
  state: ProgressState,
  action: ProgressAction
): ProgressState => {
  switch (action.type) {
    case "UPDATE_IMAGE":
      toast.dismiss();
      const toastID = toast.loading("Processing Image...", {
        position: "bottom-right",
        autoClose: false,
        hideProgressBar: false,
        closeOnClick: false,
        pauseOnHover: true,
        draggable: false,
        progress: undefined,
      });
      return {
        id: action.payload,
        status: "PROCESSING",
        notifyReference: toastID,
      };
    case "UPDATE_PROGRESS":
      return {
        status: action.payload,
      };
    default:
      throw new Error();
  }
};

export const ProgressContextProvider = (
  props: PropsWithChildren<any>
): JSX.Element => {
  const [state, dispatch] = useReducer<
    React.Reducer<ProgressState, ProgressAction>
  >(reducer, {} as ProgressState);
  const intervalRef = useRef<any>(undefined);

  useEffect(() => {
    if (intervalRef.current && state.status == "PROCESSING") {
      intervalRef.current = setInterval(() => {
        if (state.notifyReference) {
          toast.dismiss();
          toast.update(state.notifyReference, {
            type: toast.TYPE.SUCCESS,
            position: "bottom-right",
            autoClose: false,
            hideProgressBar: false,
            closeOnClick: false,
            pauseOnHover: true,
            draggable: false,
            progress: undefined,
          });
        }
      }, 15000);
    } else if (intervalRef.current && state.status !== "PROCESSING") {
      clearInterval(intervalRef.current);
    }
  }, [state]);

  return (
    <ProgressContext.Provider value={{ state, dispatch }}>
      {props.children}
    </ProgressContext.Provider>
  );
};
