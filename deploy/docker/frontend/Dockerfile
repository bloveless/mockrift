FROM node:12.12-alpine

RUN apk add --no-cache curl gnupg \
    && mkdir /app \
    && chown -R node:node /app

WORKDIR /app

USER node

RUN touch ~/.profile \
    && curl -o- -L https://yarnpkg.com/install.sh | ash

ENV PATH="/app/node_modules/.bin:$HOME/.yarn/bin:$HOME/.config/yarn/global/node_modules/.bin:$PATH"

COPY react-ui/package.json /app/package.json
RUN yarn install \
    && yarn global add react-scripts@3.0.1

CMD ["yarn", "start"]
