import React, {useRef} from "react";
import {useSubscribeNewsletter} from "./hooks/use-subscribe-newsletter";
import classnames from "classnames";

export const SubscribeFormLarge: React.FC = () => {
  const opts = {
    endingEmailSubscriptionTitle: "Interested in growing with us?",
    endingEmailSubscriptionCtaBtn: "Subscribe",
    endingEmailSubscriptionCtaLink: "#",
    endingEmailSubscriptionText: "10x Newsletter - The latest technology and business news, articles, and resources, sent straight to your inbox every month.",

  }
  const inputRef = useRef<HTMLInputElement>(null);
  const [{data, error}, subscribe] = useSubscribeNewsletter();

  return (
    <div className="container">
      <div className="row justify-content-center mb-0 mb-md-3">
        <div className="col-xl-6 col-lg-8 col-md-10 text-center">
          <h3 className="h1">{opts.endingEmailSubscriptionTitle}</h3>
        </div>
      </div>
      <div className="row justify-content-center text-center">
        <div className="col-xl-6 col-lg-7 col-md-9">
          <form onSubmit={async (e) => {
            e.preventDefault();
            await subscribe(inputRef.current?.value);
          }} className="d-md-flex mb-3 justify-content-center">
            <input
              ref={inputRef}
              type="email"
              className="mx-1 mb-2 mb-md-0 form-control form-control-lg"
              placeholder="Email Address"
            />
            <button className="mx-1 btn btn-primary-3 btn-lg" type="submit">
              {opts.endingEmailSubscriptionCtaBtn}
            </button>
          </form>

          <div
            className={classnames("alert alert-success", {"d-none": !(data && !data?.errors)})}
            role="alert"
            data-success-message=""
          >
            Thanks for your subscription!
          </div>
          <div
            className={classnames("alert alert-danger", {"d-none": !(error || data?.errors)})}
            role="alert"
            data-error-message=""
          >
            Please fill all fields correctly.
          </div>

          <div className="text-small text-muted mx-xl-6">
            {opts.endingEmailSubscriptionText}
          </div>
        </div>
      </div>
    </div>

  );
}
