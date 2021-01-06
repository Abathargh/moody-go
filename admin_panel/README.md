# admin panel

This folder contains a react.js application that interfaces with the gateway to offer a browser based gui where the gateway settings are exposed. 

## Contents

- [Requirements](#requirements)
- [Installation](#installation)
    - [As a standalone application](#As-a-standalone-application)
    - [As part of the moody front architecture](#as-part-of-the-moody-front-architecture)
- [Navigation](#navigation)
    - [Monitor](#monitor)
    - [Services & Situations](#services-&-situations)
    - [Neural](#neural)
    - [Manage Actuators](#manage-actuators)

## Requirements

Installing the app directly requires node/npm/make; you can also run it through the available pre-built docker images, or by building one yourself with the Dockerfile shipped in the repo.

## Installation

### As a standalone application

You may install the app directly:

```bash
# using make
make build
make start

# close with
make stop

# doing everything manually
cd admin_panel
npm install
npm run build

npm install -g serve
serve -l tcp:/0.0.0.0:3000 -s build/
```

Or you can just use a pre-built docker image:

```bash
docker run --name moody-adminpanel -p 3000:3000 abathargh/moody-go-adminpanel:latest
```

Docker images for this application are available for the following architectures:
- amd64
- arm32v7
- arm64

### As part of the moody front architecture

**Refer to what is written in the same section of the gateway README.**

## Navigation

The application is presented as a single page app with four different sections that group all the main features.

### Monitor

This is a simple page that contains a box for each currently active service. Each box is updated in real time with the last instance of data read from a sensor publishing on that service topic.

### Services & Situations

Every operation concerning the managing of situations and services is exposed through this section of the panel.

You can create and remove both services and situations and see a list of the currently available ones.

Services can be activated and deactivated throught the left column.

### Neural

The neural section is the center piece of the admin panel, and it's the section under which it's possible to manage the collected data and to decide wich datasets to utilize to automate the actuation process.

It's divided in four sub-sections:
 - The **Neural Monitor**, in which the current neural state is exposed. When the state is set to *stopped*, nothing happens (data is still received and shown through the monitor); when it is set to *collecting*, incoming data is added to the specified dataset, but only if the active services exactly match the dataset keys; when the state is set to *predicting*, a snapshot of the last instance of data received from each active service is acquired and forwarded to the neural service, to predict the situation that is currently happening using the selected dataset.
 - The **Situation Monitor**, where the current situation is shown and where it can also be set, according to the currently available ones. This situation is required for collecting data, and it is appended to a set of collected data when forwarding it to the dataset service.
 - The **Datasets** section, where a list of all the created datasets is shown, with the possibility to remove it and to check the keys in use; removing a dataset is final and it deletes every instance of collected data tied to that dataset.
 - The **Create Dataset** section, where it's possible to create a new dataset: the dataset keys are the currently active services.

### Manage Actuators

By default, actuators start in actuation mode. By clicking on the **Activate Server Mode** button, it's possible to switch every actuator in the local network to server mode.

When in server mode, a actuator exposes its API which can be accessed through this section.
A box appears in real time when a server mode switched actuator is detected: through these boxes, it is possible to add and delete situation to action mappings, so that it's possible to manipulate what actuator nodes should do when receiving different situations.

Every actuator has a list of actions it can perform which is encoded through integers and can be used here to tie those actions with situations.

when this process is finished, the normal actuation mode can be restored by clickin on **End Server Mode**.

