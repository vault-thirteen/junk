# Using the image from a local registry.
FROM localhost:5000/libre_office AS experiment

# Show version of important things.
RUN echo "$PATH"
RUN uname --all

# Where are we ?
RUN echo $PWD
RUN ls -l

# Prepare the test folder.
RUN mkdir "/exp"
WORKDIR "/exp"
COPY convert_source_to_pdf.sh convert_source_to_pdf.sh
RUN chmod +x convert_source_to_pdf.sh
COPY convert_pdf_to_png.sh convert_pdf_to_png.sh
RUN chmod +x convert_pdf_to_png.sh
