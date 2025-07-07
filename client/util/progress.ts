import { AxiosProgressEvent } from "axios";
import { useMemo, useState } from "react";

export interface ProgressState {
  percentage: number;
  estimatedSeconds: number;
  speed: number;
}

export interface ProgressHandle {
  state: ProgressState;
  updateProgressState: (progressEvent: AxiosProgressEvent) => void;
  reset: () => void;
  setFinished: () => void;
}

export function useProgressHandle(): ProgressHandle {
  const [progressPercentage, setProgressPercentage] = useState(0);
  const [speed, setSpeed] = useState(0);
  const [estimatedSeconds, setEstimatedSeconds] = useState(0);

  return useMemo(
    () => ({
      state: { percentage: progressPercentage, estimatedSeconds, speed },
      updateProgressState(progressEvent) {
        if (progressEvent.rate && progressEvent.rate !== speed) {
          setSpeed(progressEvent.rate);
        }

        if (
          progressEvent.progress &&
          progressEvent.progress !== progressPercentage
        ) {
          setProgressPercentage(progressEvent.progress);
        }

        if (
          progressEvent.estimated &&
          progressEvent.estimated !== estimatedSeconds
        ) {
          setEstimatedSeconds(progressEvent.estimated);
        }
      },
      reset() {
        setEstimatedSeconds(0);
        setProgressPercentage(0);
        setSpeed(0);
      },
      setFinished() {
        setEstimatedSeconds(0);
        setProgressPercentage(1);
      },
    }),

    [estimatedSeconds, progressPercentage, speed],
  );
}
