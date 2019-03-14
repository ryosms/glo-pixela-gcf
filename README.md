# glo-pixela-gcf
Update Pixela graph with google cloud functions, if glo boards is updated.

### Prepare

1. Setup GCP account to use cloud functions.
1. Setup GCloud commandline tool (should be enable `beta`).
    * If you don't setup gcloud, you can deploy function with gcp web console.

### Getting Started

1. Clone this repository

    ```bash
    $ git clone https://github.com/ryosms/glo-pixela-gcf.git
    ```
    
1. Setup a graph on [Pixela](https://pixe.la)
    1. If you don't have a pixela account, create an account.
        * <https://docs.pixe.la/#/post-user>
    1. Create a graph with `int` type and `none` selfSufficient.
        * <https://docs.pixe.la/#/post-graph>
    1. Create both webhooks `increment` and `decrement`

    * If you use JetBrains' IDE, you can create a graph on your IDE. See [pixela/README](./pixela/README.md)

1. Setup your environment variables
    1. Copy the file `sample.env.yml` to `.env.yml`
    1. Edit `.env.yml` for your environment

1. Deploy function

    ```bash
    $ cd cloudfunctions
    $ gcloud functions deploy glo-to-pixela \
        --runtime go111 \
        --entry-point GloToPixela \
        --trigger-http \
        --region <YOUR-REGION-HERE> \
        --env-vars-file ../.env.yml
    ```

    * After deploy success, the function has endpoint. ex: https://<region>-<project-id>.cloudfunctions.net/glo-to-pixela

1. Setup [Glo Board](https://app.gitkraken.com/glo/) webhook.

    1. Create a board if you have no boards.
    1. Add a webhook on `Board Settings`.
        * Name: `any`
        * Payload URL: the function's endpoint
        * Content Type: `application/json`
        * Trigger Event: only `Card`

1. Use Glo Board
    * Add Cards
    * Delete Cards
    * Archive Cards
    * and more...

1. Pixela Graph is updated!
    * Example
    
    [![](https://pixe.la/v1/users/ryosms/graphs/glo-card-status)](https://pixe.la/v1/users/ryosms/graphs/glo-card-status.html)