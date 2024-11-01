def main():
    for i in range(65535):
        try:
            print(f"{i} {chr(i)}", end=" ")
        except Exception:
            pass


if __name__ == "__main__":
    main()
