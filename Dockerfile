FROM golang:1.20-alpine   AS builder

  
WORKDIR /app
#RUN mkdir -p /run/secrets
#RUN echo "yourpassword" > /run/secrets/db-password
# to cache dependencies layer to save time for later builds if dependencies stay the same
COPY go.mod go.sum ./  
RUN go mod download
COPY . .                                                      
RUN CGO_ENABLED=0 GOOS=linux go build -o main .       

FROM scratch
# RUN mkdir -p /run/secrets
# RUN echo "yourpassword" > /run/secrets/db-password
COPY --from=builder /app/main /main
#COPY --from=builder /run/secrets/db-password /run/secrets/db-password
EXPOSE 8000
CMD [ "./main" ]


#####################################

# FROM golang:1.20-alpine   AS builder

  
# WORKDIR /app
# RUN mkdir -p /run/secrets
# RUN echo "yourpassword" > /run/secrets/db-password
# # to cache dependencies layer to save time for later builds if dependencies stay the same
# COPY go.mod go.sum ./  
# RUN go mod download
# COPY . .                                                      
# RUN CGO_ENABLED=0 GOOS=linux go build -o main .       

# FROM nginx:alpine
# # RUN mkdir -p /run/secrets
# # RUN echo "yourpassword" > /run/secrets/db-password
# COPY nginx/nginx.conf /etc/nginx/nginx.conf
# COPY nginx_certificates/ /etc/nginx/certs/ 
# COPY --from=builder /app/main /main
# COPY --from=builder /run/secrets/db-password /run/secrets/db-password
# EXPOSE 8000
# EXPOSE 443
# CMD [ "./main" ]