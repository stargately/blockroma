import React, {useRef} from "react";
import {useSubscribeNewsletter} from "./hooks/use-subscribe-newsletter";
import classnames from "classnames";

export const SubscribeForm: React.FC = () => {
  const opts = {
    subscribeTitle: "10x Newsletter",
    subscribeDesc: "The latest technology and business news, articles, and resources, sent straight to your inbox every month.",
    subscribeButton: "Subscribe",
    subscribeLink: "#",
    subscribePrivacyPolicyHtml: `Check our <a rel="noreferrer noopener nofollow" target="_blank" href="https://tianpan.co/">previous posts</a>`,

  }
  const inputRef = useRef<HTMLInputElement>(null);
  const [{data, error}, subscribe] = useSubscribeNewsletter();

  return (
    <div className="col-lg-5 col-xl-4 mt-3 mt-lg-0">
      <h5>{opts.subscribeTitle}</h5>
      <p>
        {opts.subscribeDesc}
      </p>
      <form onSubmit={async (e) => {
        e.preventDefault();
        await subscribe(inputRef.current?.value);
      }} data-form-email="" noValidate>
        <div className="form-row">
          <div className="col-12">
            <input
              ref={inputRef}

              type="email"
              className="form-control mb-2"
              placeholder="Email Address"
              name="email"
              required
            />
          </div>
          <div className="col-12">
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
            <div
              data-recaptcha=""
              data-sitekey="INSERT_YOUR_RECAPTCHA_V2_SITEKEY_HERE"
              data-size="invisible"
              data-badge="bottomleft"
            ></div>
            <button
              type="submit"
              className="btn btn-primary btn-loading btn-block"
              data-loading-text="Sending"

            >
              <img
                className="icon"
                src="/assets/img/icons/theme/code/loading.svg"
                alt="loading icon"
                data-inject-svg=""
              />
              <span>Subscribe</span>
            </button>
          </div>
        </div>
      </form>
      <small className="text-muted form-text" dangerouslySetInnerHTML={{__html: opts.subscribePrivacyPolicyHtml}}>
      </small>
    </div>

  );
}
