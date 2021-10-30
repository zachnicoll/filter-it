import { toastLoading } from "common/toast";
import { ReactText } from "react";
import { toast } from "react-toastify";
import produce from "immer";
import { Progress } from "api/types";

export type ProgressAction =
  | { type: "UPDATE_IMAGE"; payload: { id: string } }
  | { type: "UPDATE_PROGRESS"; payload: { status: Progress } };

export interface ProgressContextState {
  id?: string;
  status?: Progress;
  notifyReference?: ReactText;
  imageLocation?: string;
}

const progressReducer = produce(
  (state: ProgressContextState, action: ProgressAction) => {
    switch (action.type) {
      case "UPDATE_IMAGE":
        toast.dismiss();
        const toastID = toastLoading("Processing Image...");

        Object.assign(state, {
          id: action.payload.id,
          status: Progress.PROCESSING,
          notifyReference: toastID,
        });
        break;

      case "UPDATE_PROGRESS":
        state.status = action.payload.status;
        break;
    }
  }
);

export default progressReducer;
