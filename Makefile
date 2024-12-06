
all:
	cd client && make && cp test ..
	cd server && make && cp compute ..